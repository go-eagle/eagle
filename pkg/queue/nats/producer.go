package nats

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// Producer define a nats producer
type Producer struct {
	addr      string
	conn      *nats.Conn
	connClose chan bool
	quit      chan struct{}
}

// NewProducer create a producer
func NewProducer(addr string) *Producer {
	p := &Producer{
		addr:      addr,
		connClose: make(chan bool),
		quit:      make(chan struct{}),
	}
	if err := p.Start(); err != nil {
		log.Println("nats start producer err: ", err)
	}
	return p
}

// Start .
func (p *Producer) Start() error {
	if err := p.Run(); err != nil {
		return err
	}

	log.Println("nats producer connected and running!")

	go p.ReConnect()
	return nil
}

// Stop .
func (p *Producer) Stop() {
	close(p.quit)
	if !p.conn.IsClosed() {
		p.conn.Close()
	}
}

// Run .
func (p *Producer) Run() error {
	var err error
	opts := nats.Options{
		MaxReconnect: -1,
		ClosedCB: func(conn *nats.Conn) {
			p.connClose <- true
			log.Println("nats producer - connection closed cb")
		},
		DisconnectedErrCB: func(conn *nats.Conn, err error) {
			log.Println("nats producer - connection disconnected err cb")
		},
		ReconnectedCB: func(conn *nats.Conn) {
			log.Println("nats producer - connection reconnected cb")
		},
		AsyncErrorCB: func(conn *nats.Conn, sub *nats.Subscription, err error) {
			log.Println("nats producer - connection async err cb")
		},
	}
	p.conn, err = opts.Connect()
	return err
}

// ReConnect .
func (p *Producer) ReConnect() {
	for {
		select {
		case closed := <-p.connClose:
			if closed {
				log.Println("nats producer - connection closed")
			}
		case <-p.quit:
			return
		}

		if !p.conn.IsClosed() {
			p.conn.Close()
		}

	quit:
		for {
			select {
			case <-p.quit:
				return
			default:
				log.Println("nats producer - reconnect")

				if err := p.Run(); err != nil {
					log.Println("nats producer - failCheck: ", err)

					// sleep 5s reconnect
					time.Sleep(time.Second * 5)
					continue
				}
				log.Println("nats producer connected and running!")
				break quit
			}
		}
	}
}

// Publish push data to queue
func (p *Producer) Publish(topic string, data interface{}) error {
	encodeConn, err := nats.NewEncodedConn(p.conn, nats.JSON_ENCODER)
	if err != nil {
		return err
	}
	return encodeConn.Publish(topic, data)
}
