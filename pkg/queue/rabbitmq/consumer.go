package rabbitmq

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

// Consumer define consumer for rabbitmq
type Consumer struct {
	addr          string
	conn          *amqp.Connection
	channel       *amqp.Channel
	connNotify    chan *amqp.Error
	channelNotify chan *amqp.Error
	quit          chan struct{}
	exchange      string
	routingKey    string
	queueName     string
	consumerTag   string
	autoDelete    bool                    // 是否自动删除
	handler       func(body []byte) error // 业务自定义消费函数
}

// NewConsumer instance a consumer
func NewConsumer(addr, exchange, queueName string, autoDelete bool, handler func(body []byte) error) *Consumer {
	return &Consumer{
		addr:        addr,
		exchange:    exchange,
		routingKey:  "",
		queueName:   queueName,
		consumerTag: "consumer",
		autoDelete:  autoDelete,
		handler:     handler,
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
			log.Println("rabbitmq consumer - channel cancel failed: ", err)
		}

		if err := c.conn.Close(); err != nil {
			log.Println("rabbitmq consumer - connection close failed: ", err)
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

	var delivery <-chan amqp.Delivery
	// NOTE: autoAck param
	delivery, err = c.channel.Consume(c.queueName, c.consumerTag, true, false, false, false, nil)
	if err != nil {
		_ = c.channel.Close()
		_ = c.conn.Close()
		return err
	}

	go c.Handle(delivery)

	c.connNotify = c.conn.NotifyClose(make(chan *amqp.Error))
	c.channelNotify = c.channel.NotifyClose(make(chan *amqp.Error))

	return nil
}

// Handle handle data
func (c *Consumer) Handle(delivery <-chan amqp.Delivery) {
	for d := range delivery {
		log.Printf("Consumer received a message: %s in queue: %s", d.Body, c.queueName)
		log.Printf("got %dB delivery: [%v] %q", len(d.Body), d.DeliveryTag, d.Body)
		go func(delivery amqp.Delivery) {
			if err := c.handler(delivery.Body); err == nil {
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
	log.Println("handle: async deliveries channel closed")
}

// ReConnect .
func (c *Consumer) ReConnect() {
	for {
		select {
		case err := <-c.connNotify:
			if err != nil {
				log.Fatalf("rabbitmq consumer - connection NotifyClose: %+v", err)
			}
		case err := <-c.channelNotify:
			if err != nil {
				log.Fatalf("rabbitmq consumer - channel NotifyClose: %+v", err)
			}
		case <-c.quit:
			return
		}

		// backstop
		if !c.conn.IsClosed() {
			// 关闭 SubMsg message delivery
			if err := c.channel.Cancel(c.consumerTag, true); err != nil {
				log.Fatalf("rabbitmq consumer - channel cancel failed: %+v", err)
			}
			if err := c.conn.Close(); err != nil {
				log.Fatalf("rabbitmq consumer - conn cancel failed: %+v", err)
			}
		}

		// IMPORTANT: 必须清空 Notify，否则死连接不会释放
		for err := range c.channelNotify {
			println(err)
		}
		for err := range c.connNotify {
			println(err)
		}

	quit:
		for {
			select {
			case <-c.quit:
				return
			default:
				log.Println("rabbitmq consumer - reconnect")

				if err := c.Run(); err != nil {
					log.Printf("rabbitmq consumer - failCheck: %+v", err)

					// sleep 5s reconnect
					time.Sleep(time.Second * 5)
					continue
				}

				break quit
			}
		}
	}
}
