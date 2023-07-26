package rabbitmq

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/go-eagle/eagle/pkg/queue/rabbitmq/options"

	"github.com/cenkalti/backoff/v4"
	"github.com/rabbitmq/amqp091-go"

	"github.com/go-eagle/eagle/pkg/log"
)

type Connection struct {
	options   *options.ConnectionOptions
	conn      *amqp091.Connection
	connected chan struct{}
	closing   chan struct{}
	closed    int32

	logger  log.Logger
	mu      sync.RWMutex
	backoff backoff.BackOff
}

func NewConnection(opts *options.ConnectionOptions, logger log.Logger) (*Connection, error) {
	conn := &Connection{
		options:   opts,
		connected: make(chan struct{}),
		closing:   make(chan struct{}),
		logger:    logger,
		backoff:   backoff.NewExponentialBackOff(),
	}
	if err := conn.connect(); err != nil {
		return nil, err
	}

	go conn.watch()

	return conn, nil
}

func (c *Connection) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		amqpConn *amqp091.Connection
		err      error
	)
	if c.options.Config != nil {
		amqpConn, err = amqp091.DialConfig(c.options.URI, c.conn.Config)
	} else {
		amqpConn, err = amqp091.Dial(c.options.URI)
	}
	if err != nil {
		return err
	}

	c.conn = amqpConn

	close(c.connected)

	return nil
}

func (c *Connection) Conn() *amqp091.Connection {
	return c.conn
}

func (c *Connection) watch() {
	for {
		<-c.connected

		select {
		case <-c.closing:
			c.logger.Info("rabbitmq: stop watch connection")
			c.connected = make(chan struct{})
			return
		case err := <-c.conn.NotifyClose(make(chan *amqp091.Error)):
			// if RabbitMQ server is shutdown right now, will receive err:
			// Exception (320) Reason: "CONNECTION_FORCED - broker forced connection closure with reason 'shutdown'"
			c.logger.Errorf("rabbitmq: received close notification from AMQP, reconnecting, err: %v", err)
			c.connected = make(chan struct{})
			c.reconnect()
		}
	}

}

func (c *Connection) reconnect() {
	reconnect := func() error {
		c.logger.Info("rabbitmq: reconnecting connection")
		// if RabbitMQ server is shutdown, return err: Exception (501) Reason: "EOF"
		// or err: dial tcp [::1]:5672: connect: connection refused
		err := c.connect()
		if err == nil {
			c.logger.Info("rabbitmq: reconnected connection successfully")
			return nil
		}
		c.logger.Errorf("rabbitmq: reconnect failed, retrying, err: %v", err)
		if c.IsClosed() {
			return backoff.Permanent(fmt.Errorf("rabbitmq: connection is closed, err: %v", err))
		}

		return err
	}

	err := backoff.Retry(reconnect, c.backoff)
	if err != nil {
		c.logger.Errorf("rabbitmq: reconnect failed: %v", err)
	}
}

func (c *Connection) IsConnected() bool {
	select {
	case <-c.connected:
		return true
	default:
		return false
	}
}

func (c *Connection) IsClosed() bool {
	return atomic.LoadInt32(&c.closed) == 1
}

func (c *Connection) Close() error {
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return nil
	}
	close(c.closing)
	c.logger.Info("rabbitmq: closing AMQP connection")
	if err := c.conn.Close(); err != nil {
		c.logger.Errorf("rabbitmq: close connection failed, err: %v", err)
		return err
	}
	c.logger.Info("rabbitmq: closed AMQP connection")
	return nil
}
