package rabbitmq

import (
	"context"
	"errors"
	"sync"

	"github.com/cenkalti/backoff/v4"

	"github.com/rabbitmq/amqp091-go"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/queue/rabbitmq/options"
)

// Action define action
type Action int

// Handler define handler for rabbitmq
type Handler func(ctx context.Context, msg amqp091.Delivery) (action Action)

const (
	// Ack default ack this msg after you have successfully processed this delivery.
	Ack Action = iota
	// NackDiscard the message will be dropped or delivered to a server configured dead-letter queue.
	NackDiscard
	// NackRequeue deliver this message to a different consumer.
	NackRequeue
)

// Consumer define consumer for rabbitmq
type Consumer struct {
	channel *Channel
	options *options.ConsumerOptions
	logger  log.Logger

	handlerWG sync.WaitGroup
	watchWG   sync.WaitGroup
	closing   chan struct{}
}

// NewConsumer instance a consumer
func NewConsumer(conf *Config, logger log.Logger) (*Consumer, error) {
	conn, err := NewConnection(conf.Connection, logger)
	if err != nil {
		return nil, err
	}
	ch, err := NewChannel(conn, conf, logger)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		channel:   ch,
		logger:    logger,
		handlerWG: sync.WaitGroup{},
		watchWG:   sync.WaitGroup{},
		closing:   make(chan struct{}),
	}, nil
}

func (c *Consumer) Consume(ctx context.Context, handler Handler, opts ...options.ConsumerOption) error {
	c.logger.Info("rabbitmq: Consumer start consuming")

	consumerOptions := options.NewConsumerOptions(opts...)
	c.options = consumerOptions

	// use one or multiple goroutines to handle deliveries
	if err := c.parallelHandle(ctx, handler); err != nil {
		c.logger.Errorf("rabbitmq: Consumer start parallelHandle error: %v", err)
		return err
	}

	go c.watch(handler)

	c.handlerWG.Wait()
	c.watchWG.Wait()

	return nil
}

func (c *Consumer) parallelHandle(ctx context.Context, handler Handler) error {
	if !c.channel.IsConnected() {
		return errors.New("rabbitmq: channel is not connected")
	}
	c.logger.Info("rabbitmq: consumer start parallelHandle")
	if err := c.channel.Qos(c.options.QOSPrefetchCount, c.options.QOSPrefetchSize, c.options.QOSGlobal); err != nil {
		c.logger.Errorf("rabbitmq: consumer set qos error: %v", err)
		return err
	}

	queue := c.channel.opts.Queue.Name
	if c.options.ConsumerQueue != "" {
		queue = c.options.ConsumerQueue
	}

	messages, err := c.channel.Consume(queue, c.options.ConsumerName, c.options.ConsumerAutoAck,
		c.options.ConsumerExclusive, c.options.ConsumerNoLocal, c.options.ConsumerNoWait,
		c.options.ConsumerArgs)
	if err != nil {
		c.logger.Errorf("rabbitmq: consumer start consume error: %v", err)
		return err
	}

	for i := 0; i < c.options.Concurrency; i++ {
		c.handlerWG.Add(1)
		go c.handle(ctx, messages, handler)
	}

	return nil
}

// Handle data
func (c *Consumer) handle(ctx context.Context, messages <-chan amqp091.Delivery, handler Handler) {
	defer c.handlerWG.Done()

	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				return
			}
			if c.options.ConsumerAutoAck {
				handler(ctx, msg)
				return
			}
			action := handler(ctx, msg)
			switch action {
			case Ack:
				if err := msg.Ack(false); err != nil {
					c.logger.Errorf("rabbitmq: consumer ack error: %v", err)
				}
			case NackDiscard:
				if err := msg.Nack(false, false); err != nil {
					c.logger.Errorf("rabbitmq: consumer nack error: %v", err)
				}
			case NackRequeue:
				if err := msg.Nack(false, true); err != nil {
					c.logger.Errorf("rabbitmq: consumer nack error: %v", err)
				}
			}
		case <-c.closing:
			c.logger.Info("rabbitmq: consumer has closed")
			return
		}
	}
}

// watch .
func (c *Consumer) watch(handler Handler) {
	c.watchWG.Add(1)
	defer func() {
		c.watchWG.Done()
	}()

	for {
		select {
		case err := <-c.channel.notifyReconnected:
			c.logger.Errorf("rabbitmq: consumer begin to retry after receiving reconnect notify, error: %v", err)
			parallelHandleFunc := func() error {
				err = c.parallelHandle(context.Background(), handler)
				if err != nil {
					c.logger.Errorf("rabbitmq: consumer retry parallelHandle error: %v", err)
					return err
				}
				return nil
			}

			err = backoff.Retry(parallelHandleFunc, backoff.NewExponentialBackOff())
			if err != nil {
				c.logger.Errorf("rabbitmq: consumer watch retry error: %v", err)
			} else {
				c.logger.Info("rabbitmq: consumer watch retry successfully")
			}
		case <-c.closing:
			c.logger.Info("rabbitmq: watch consumer has closed")
			return
		}
	}
}

// IsClosed check consumer is closed
func (c *Consumer) IsClosed() bool {
	select {
	case <-c.closing:
		return true
	default:
		return false
	}
}

// Close consumer
func (c *Consumer) Close() error {
	c.logger.Info("rabbitmq: Consumer is closing")
	close(c.closing)
	if err := c.channel.Close(); err != nil {
		return err
	}
	c.logger.Info("rabbitmq: Consumer closed successfully")
	return nil
}
