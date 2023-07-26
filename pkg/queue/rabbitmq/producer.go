package rabbitmq

import (
	"context"
	"errors"
	"sync"

	"github.com/cenkalti/backoff/v4"
	"github.com/go-eagle/eagle/pkg/queue/rabbitmq/options"
	"github.com/rabbitmq/amqp091-go"

	"github.com/go-eagle/eagle/pkg/log"
)

// Producer define struct for producer
type Producer struct {
	channel    *Channel
	exchange   string
	routingKey string
	mu         sync.Mutex
}

// NewProducer create a producer
func NewProducer(conf *Config, logger log.Logger) (*Producer, error) {
	conn, err := NewConnection(conf.Connection, logger)
	if err != nil {
		return nil, err
	}
	ch, err := NewChannel(conn, conf, logger)
	if err != nil {
		return nil, err
	}

	if ch.opts.Exchange.Name == "" || ch.opts.Bind.RoutingKey == "" {
		return nil, errors.New("exchange name or routing key is empty")
	}

	p := &Producer{
		channel:    ch,
		exchange:   ch.opts.Exchange.Name,
		routingKey: ch.opts.Bind.RoutingKey,
	}

	return p, nil
}

// Publish push data to queue
func (p *Producer) Publish(ctx context.Context, message []byte, opts ...options.PublishOption) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	msgOptions := options.NewPublishOption(opts...)
	msg := amqp091.Publishing{
		Headers:         msgOptions.Headers,
		ContentType:     msgOptions.ContentType,
		ContentEncoding: msgOptions.ContentEncoding,
		DeliveryMode:    msgOptions.DeliveryMode,
		Priority:        msgOptions.Priority,
		CorrelationId:   msgOptions.CorrelationId,
		ReplyTo:         msgOptions.ReplyTo,
		Expiration:      msgOptions.Expiration,
		MessageId:       msgOptions.MessageId,
		Timestamp:       msgOptions.Timestamp,
		Type:            msgOptions.Type,
		UserId:          msgOptions.UserId,
		AppId:           msgOptions.AppId,
		Body:            message,
	}

	// if set exchange and routing key in publish option, override the default value
	exchange := p.exchange
	if msgOptions.MsgExchange != "" {
		exchange = msgOptions.MsgExchange
	}
	routingKey := p.routingKey
	if msgOptions.MsgRoutingKey != "" {
		routingKey = msgOptions.MsgRoutingKey
	}

	publishFunc := func() error {
		err := p.channel.Publish(ctx, exchange, routingKey, msgOptions.Mandatory, msgOptions.Immediate, msg)
		return err
	}

	expBackoff := backoff.NewExponentialBackOff()
	retry := backoff.WithMaxRetries(expBackoff, msgOptions.MaxRetry)
	err := backoff.Retry(publishFunc, retry)
	if err != nil {
		return err
	}

	return nil
}

// Close connection
func (p *Producer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.channel.Close()
}
