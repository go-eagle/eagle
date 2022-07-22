package rabbitmq

import (
	"context"
	"time"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/streadway/amqp"
)

type Handler func(ctx context.Context, body []byte) error

// Consumer define consumer for rabbitmq
type Consumer struct {
	addr          string
	conn          *amqp.Connection
	channel       *amqp.Channel
	connNotify    chan *amqp.Error
	channelNotify chan *amqp.Error
	quit          chan struct{}
	exchange      string
	queueName     string
	consumerTag   string
	autoDelete    bool // 是否自动删除
}

// NewConsumer instance a consumer
func NewConsumer(addr, exchange string, autoDelete bool) *Consumer {
	return &Consumer{
		addr:        addr,
		exchange:    exchange,
		consumerTag: "consumer",
		autoDelete:  autoDelete,
		quit:        make(chan struct{}),
	}
}

// Start start a service
func (c *Consumer) Start() error {
	if err := c.Run(); err != nil {
		return err
	}

	go c.ReConnect()

	return nil
}

// Stop a consumer
func (c *Consumer) Stop() {
	close(c.quit)

	if !c.conn.IsClosed() {
		// 关闭 SubMsg message delivery
		if err := c.channel.Cancel(c.consumerTag, true); err != nil {
			log.Errorf("rabbitmq consumer - channel cancel failed: ", err)
		}

		if err := c.conn.Close(); err != nil {
			log.Errorf("rabbitmq consumer - connection close failed: ", err)
		}
	}
}

// Run .
func (c *Consumer) Run() error {
	var err error
	if c.conn, err = OpenConnection(c.addr); err != nil {
		return err
	}

	if c.channel, err = NewChannel(c.conn).Create(); err != nil {
		_ = c.conn.Close()
		return err
	}

	// bind queue in rabbitmq admin ui
	//if _, err = c.channel.QueueDeclare(c.queueName, true, c.autoDelete, false, false, nil); err != nil {
	//	_ = c.channel.Close()
	//	_ = c.conn.Close()
	//	return err
	//}
	//
	//if err = c.channel.QueueBind(c.queueName, c.routingKey, c.exchange, false, nil); err != nil {
	//	_ = c.channel.Close()
	//	_ = c.conn.Close()
	//	return err
	//}

	return nil
}

func (c *Consumer) Consume(ctx context.Context, queueName string, handler Handler) {
	var (
		err      error
		delivery <-chan amqp.Delivery
	)

	c.connNotify = c.conn.NotifyClose(make(chan *amqp.Error))
	c.channelNotify = c.channel.NotifyClose(make(chan *amqp.Error))

	for {
		// NOTE: autoAck param
		delivery, err = c.channel.Consume(
			queueName,
			c.consumerTag,
			true,
			false,
			false,
			false,
			nil)
		if err != nil {
			log.Errorf("Consumer channel Consume err: %#v", err)
			time.Sleep(5 * time.Second)
		}

		time.Sleep(5 * time.Second)

		go c.Handle(delivery, handler)
	}

}

// Handle handle data
func (c *Consumer) Handle(delivery <-chan amqp.Delivery, handler Handler) {
	ctx := context.Background()
	for d := range delivery {
		log.Infof("Consumer received a message: %s in queue: %s", d.Body, c.queueName)
		log.Infof("got %dB delivery: [%v] %q", len(d.Body), d.DeliveryTag, d.Body)
		go func(delivery amqp.Delivery) {
			if err := handler(ctx, delivery.Body); err == nil {
				// NOTE: 假如现在有 10 条消息，它们都是并发处理的，如果第 10 条消息最先处理完毕，
				// 那么前 9 条消息都会被 delivery.Ack(true) 给确认掉。后续 9 条消息处理完毕时，
				// 再执行 delivery.Ack(true)，显然就会导致消息重复确认
				// 报 406 PRECONDITION_FAILED 错误， 所以这里为 false
				_ = delivery.Ack(false)
			} else {
				// 重新入队，否则未确认的消息会持续占用内存
				_ = delivery.Reject(true)
			}
		}(d)
	}
	log.Infof("handle: async deliveries channel closed")
}

// ReConnect .
func (c *Consumer) ReConnect() {
	for {
		select {
		case err := <-c.connNotify:
			if err != nil {
				log.Errorf("rabbitmq consumer - connection NotifyClose: %+v", err)
			}
		case err := <-c.channelNotify:
			if err != nil {
				log.Errorf("rabbitmq consumer - channel NotifyClose: %+v", err)
			}
		case <-c.quit:
			return
		}

		// backstop
		if !c.conn.IsClosed() {
			// 关闭 SubMsg message delivery
			if err := c.channel.Cancel(c.consumerTag, true); err != nil {
				log.Errorf("rabbitmq consumer - channel cancel failed: %+v", err)
			}
			if err := c.conn.Close(); err != nil {
				log.Errorf("rabbitmq consumer - conn cancel failed: %+v", err)
			}
		}

		// IMPORTANT: 必须清空 Notify，否则死连接不会释放
		for err := range c.channelNotify {
			log.Errorf("rabbitmq consumer - channelNotify err: %+v", err)
		}
		for err := range c.connNotify {
			log.Errorf("rabbitmq consumer - connNotify err: %+v", err)
		}

	quit:
		for {
			select {
			case <-c.quit:
				return
			default:
				log.Infof("rabbitmq consumer - reconnect")

				if err := c.Run(); err != nil {
					log.Infof("rabbitmq consumer - failCheck: %+v", err)

					// sleep 5s reconnect
					time.Sleep(time.Second * 5)
					continue
				}

				break quit
			}
		}
	}
}
