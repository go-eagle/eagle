package nats

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// Consumer define a nats consumer
type Consumer struct {
	addr      string
	conn      *nats.Conn
	subscribe *nats.Subscription
	connClose chan bool
	quit      chan struct{}
}

// NewConsumer create consumer
func NewConsumer(addr string) *Consumer {
	c := &Consumer{
		addr:      addr,
		connClose: make(chan bool),
		quit:      make(chan struct{}),
	}
	if err := c.Start(); err != nil {
		log.Println("nats start consumer err: ", err)
	}
	return c
}

// Start .
func (c *Consumer) Start() error {
	if err := c.Run(); err != nil {
		return err
	}

	log.Println("nats consumer connected and running!")

	go c.ReConnect()
	return nil
}

// Stop .
func (c *Consumer) Stop() {
	close(c.quit)
	if !c.conn.IsClosed() {
		_ = c.subscribe.Unsubscribe()
		c.conn.Close()
	}
}

// Run .
func (c *Consumer) Run() error {
	var err error
	opts := nats.Options{
		MaxReconnect: -1,
		ClosedCB: func(conn *nats.Conn) {
			c.connClose <- true
			log.Println("nats consumer - connection closed cb")
		},
		DisconnectedErrCB: func(conn *nats.Conn, err error) {
			log.Println("nats consumer - connection disconnected err cb")
		},
		ReconnectedCB: func(conn *nats.Conn) {
			log.Println("nats consumer - connection reconnected cb")
		},
		AsyncErrorCB: func(conn *nats.Conn, sub *nats.Subscription, err error) {
			log.Println("nats consumer - connection async err cb")
		},
	}
	c.conn, err = opts.Connect()
	return err
}

// ReConnect .
func (c *Consumer) ReConnect() {
	for {
		select {
		case closed := <-c.connClose:
			if closed {
				log.Println("nats consumer - connection closed")
			}
		case <-c.quit:
			return
		}

		if !c.conn.IsClosed() {
			c.conn.Close()
		}

	quit:
		for {
			select {
			case <-c.quit:
				return
			default:
				log.Println("nats consumer - reconnect")

				if err := c.Run(); err != nil {
					log.Println("nats consumer - failCheck: ", err)

					// sleep 5s reconnect
					time.Sleep(time.Second * 5)
					continue
				}
				log.Println("nats consumer connected and running!")
				break quit
			}
		}
	}
}

// Consume consume data from nats queue
func (c *Consumer) Consume(topic string, handler interface{}) error {
	encodeConn, err := nats.NewEncodedConn(c.conn, nats.JSON_ENCODER)
	if err != nil {
		return err
	}

	if c.subscribe != nil {
		_ = c.subscribe.Unsubscribe()
		c.subscribe = nil
	}

	c.subscribe, err = encodeConn.Subscribe(topic, handler)
	if err != nil {
		return err
	}
	_ = encodeConn.Flush()

	return nil
}
