package rabbitmq

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// Producer define struct for rabbitmq
type Producer struct {
	addr          string
	conn          *amqp.Connection
	channel       *amqp.Channel
	routingKey    string
	exchange      string
	connNotify    chan *amqp.Error
	channelNotify chan *amqp.Error
	quit          chan struct{}
}

// NewProducer create a producer
func NewProducer(addr, exchange string) *Producer {
	p := &Producer{
		addr:     addr,
		exchange: exchange,
		quit:     make(chan struct{}),
	}

	return p
}

// Start start a producer
func (p *Producer) Start() error {
	if err := p.Run(); err != nil {
		return err
	}
	go p.ReConnect()

	return nil
}

// Stop .
func (p *Producer) Stop() {
	close(p.quit)

	if !p.conn.IsClosed() {
		if err := p.conn.Close(); err != nil {
			log.Println("rabbitmq producer - connection close failed: ", err)
		}
	}
}

// Run .
func (p *Producer) Run() error {
	var err error
	if p.conn, err = OpenConnection(p.addr); err != nil {
		return err
	}

	if p.channel, err = NewChannel(p.conn).Create(); err != nil {
		_ = p.conn.Close()
		return err
	}

	p.connNotify = p.conn.NotifyClose(make(chan *amqp.Error))
	p.channelNotify = p.channel.NotifyClose(make(chan *amqp.Error))

	return err
}

// ReConnect .
func (p *Producer) ReConnect() {
	for {
		select {
		case err := <-p.connNotify:
			if err != nil {
				log.Println("rabbitmq producer - connection NotifyClose: ", err)
			}
		case err := <-p.channelNotify:
			if err != nil {
				log.Println("rabbitmq producer - channel NotifyClose: ", err)
			}
		case <-p.quit:
			return
		}

		// backstop
		if !p.conn.IsClosed() {
			if err := p.conn.Close(); err != nil {
				log.Println("rabbitmq producer - connection close failed: ", err)
			}
		}

		// IMPORTANT: 必须清空 Notify，否则死连接不会释放
		for err := range p.channelNotify {
			log.Println(err)
		}
		for err := range p.connNotify {
			log.Println(err)
		}

	quit:
		for {
			select {
			case <-p.quit:
				return
			default:
				log.Println("rabbitmq producer - reconnect")

				if err := p.Run(); err != nil {
					log.Println("rabbitmq producer - failCheck: ", err)

					// sleep 5s reconnect
					time.Sleep(time.Second * 5)
					continue
				}

				break quit
			}
		}
	}
}

// Publish push data to queue
func (p *Producer) Publish(routingKey, message string) error {
	return p.channel.Publish(
		p.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			MessageId:    uuid.New().String(),
			Type:         "",
			Body:         []byte(message),
			Timestamp:    time.Now(),
		})
}
