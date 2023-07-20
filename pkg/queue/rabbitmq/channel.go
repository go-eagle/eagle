package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"

	"github.com/go-eagle/eagle/pkg/log"
)

// Channel data channel
type Channel struct {
	conn              *Connection
	ch                *amqp091.Channel
	opts              ConnectionOptions
	connected         chan struct{}
	closing           chan struct{}
	notifyReconnected chan error

	once   sync.Once
	logger log.Logger
}

// NewChannel instance a channel
func NewChannel(conn *Connection, logger log.Logger) (*Channel, error) {
	ch := &Channel{
		conn:              conn,
		connected:         make(chan struct{}),
		closing:           make(chan struct{}),
		notifyReconnected: make(chan error, 1),
		once:              sync.Once{},
		logger:            logger,
	}
	if err := ch.connect(); err != nil {
		return nil, err
	}

	go ch.watch()

	return ch, nil
}

// Create create a channel
func (c *Channel) create() (*amqp091.Channel, error) {
	return c.conn.conn.Channel()
}

// connect connect a channel
func (c *Channel) connect() error {
	channel, err := c.create()
	if err != nil {
		return err
	}
	c.ch = channel
	c.initDeclare()
	close(c.connected)
	return nil
}

// initDeclare declare a channel
func (c *Channel) initDeclare() error {
	var err error
	c.once.Do(func() {
		exchange := ""
		err = c.ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
		if err != nil {
			return
		}
		_, err = c.ch.QueueDeclare("queue", true, false, false, false, nil)
		if err != nil {
			return
		}
		err = c.ch.QueueBind("queue", "routing_key", exchange, false, nil)
	})
	return err
}

// watch monitor a channel
func (c *Channel) watch() {
	for {
		<-c.connected

		select {
		case err := <-c.ch.NotifyClose(make(chan *amqp091.Error)):
			c.logger.Errorf("rabbitmq: channel closed, err: %s", err)
			c.connected = make(chan struct{})
			// use select to avoid blocking
			select {
			case c.notifyReconnected <- err:
			default:
			}
			c.reconnect()
		case err := <-c.ch.NotifyCancel(make(chan string)):
			c.logger.Errorf("rabbitmq: channel cancelled, err: %s", err)
			c.connected = make(chan struct{})
			// use select to avoid blocking
			select {
			case c.notifyReconnected <- errors.New(err):
			default:
			}
			c.reconnect()
		case <-c.closing:
			c.logger.Info("rabbitmq: stopping watch channel")
			return
		}
	}
}

// reconnect if channel is closed, reconnect
func (c *Channel) reconnect() {
	reconnect := func() error {
		// if channel is closed, return error to stop retry
		if c.conn.IsConnected() == false {
			return fmt.Errorf("rabbitmq: connection is not connected")
		}
		if err := c.connect(); err != nil {
			return err
		}
		c.logger.Info("rabbitmq: channel reconnected")
		return nil
	}

	err := backoff.Retry(reconnect, backoff.NewExponentialBackOff())
	if err != nil {
		c.logger.Errorf("rabbitmq: channel reconnect error: %+v", err)
	}
}

// IsConnected check channel is connected
func (c *Channel) IsConnected() bool {
	select {
	case <-c.connected:
		return true
	default:
		return false
	}
}

// Publish a message
func (c *Channel) Publish(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp091.Publishing) error {
	select {
	case <-c.connected:
		return c.ch.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg)
	case <-time.After(time.Second * 5):
		return fmt.Errorf("rabbitmq: Publish msg is timeout: %+v", time.Second*5)
	}
}

// Consumer consume a message
func (c *Channel) Consumer(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp091.Table) (<-chan amqp091.Delivery, error) {
	select {
	case <-c.connected:
		return c.ch.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
	case <-time.After(time.Second * 5):
		return nil, fmt.Errorf("rabbitmq: Consumer msg is timeout: %+v", time.Second*5)
	}
}

// IsClosed check channel is closed
func (c *Channel) IsClosed() bool {
	if c.ch != nil {
		return c.ch.IsClosed()
	}
	return true
}

// Close close a channel
func (c *Channel) Close() error {
	close(c.closing)
	if c.ch != nil {
		if err := c.ch.Close(); err != nil {
			c.logger.Errorf("rabbitmq: close channel error: %+v", err)
			return err
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.logger.Errorf("rabbitmq: close connection error: %+v", err)
			return err
		}
	}

	return nil
}
