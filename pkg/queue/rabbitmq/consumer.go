package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	queueName     string
	consumerTag   string
	connNotify    chan *amqp.Error
	channelNotify chan *amqp.Error
	autoDelete    bool // 是否自动删除
	quit          chan struct{}
	isSync        bool                    // 是否同步消费
	handler       func(body []byte) error // 业务自定义消费函数
}

func NewConsumer(conn *amqp.Connection, channel *amqp.Channel, queueName string) *Consumer {
	return &Consumer{
		conn:      conn,
		channel:   channel,
		queueName: queueName,
	}
}

func (c *Consumer) Consume(isSync bool, handler func(body []byte) error) error {
	c.isSync = isSync
	c.handler = handler

	if err := c.Run(); err != nil {
		return err
	}

	go c.ReConnect()

	return nil
}

func (c *Consumer) Run() error {
	var delivery <-chan amqp.Delivery
	delivery, err := c.channel.Consume(
		c.queueName,
		c.consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.channel.Close()
		c.conn.Close()
		return fmt.Errorf("queue consume: %s", err)
	}

	if c.isSync {
		go c.syncHandle(delivery)
	} else {
		go c.asyncHandle(delivery)
	}

	c.connNotify = c.conn.NotifyClose(make(chan *amqp.Error))
	c.channelNotify = c.channel.NotifyClose(make(chan *amqp.Error))

	return nil
}

// 同步处理方式
func (c *Consumer) syncHandle(delivery <-chan amqp.Delivery) {
	for d := range delivery {
		log.Printf("Consumer received a message: %s in queue: %s", d.Body, c.queueName)
		log.Printf("got %dB delivery: [%v] %q", len(d.Body), d.DeliveryTag, d.Body)
		if err := c.handler(d.Body); err != nil {
			d.Ack(true)
		} else {
			// 重新入队，否则未确认的消息会持续占用内存
			d.Reject(true)
		}
	}
	log.Println("handle: sync deliveries channel closed")
}

// 异步处理方式
func (c *Consumer) asyncHandle(delivery <-chan amqp.Delivery) {
	for d := range delivery {
		log.Printf("Consumer received a message: %s in queue: %s", d.Body, c.queueName)
		log.Printf("got %dB delivery: [%v] %q", len(d.Body), d.DeliveryTag, d.Body)
		go func(delivery amqp.Delivery) {
			if err := c.handler(delivery.Body); err != nil {
				// NOTE: 假如现在有 10 条消息，它们都是并发处理的，如果第 10 条消息最先处理完毕，
				// 那么前 9 条消息都会被 delivery.Ack(true) 给确认掉。后续 9 条消息处理完毕时，
				// 再执行 delivery.Ack(true)，显然就会导致消息重复确认
				// 报 406 PRECONDITION_FAILED 错误， 所以这里为 false
				delivery.Ack(false)
			} else {
				// 重新入队，否则未确认的消息会持续占用内存
				delivery.Reject(true)
			}
		}(d)
	}
	log.Println("handle: async deliveries channel closed")
}

func (c *Consumer) ReConnect() {
	for {
		select {
		case err := <-c.connNotify:
			if err != nil {
				log.Fatalf("rabbitmq consumer - connection NotifyClose: ", err)
			}
		case err := <-c.channelNotify:
			if err != nil {
				log.Fatalf("rabbitmq consumer - channel NotifyClose: ", err)
			}
		case <-c.quit:
			return
		}

		// backstop
		if !c.conn.IsClosed() {
			// 关闭 SubMsg message delivery
			if err := c.channel.Cancel(c.consumerTag, true); err != nil {
				log.Fatalf("rabbitmq consumer - channel cancel failed: ", err)
			}
			if err := c.conn.Close(); err != nil {
				log.Fatalf("rabbitmq consumer - conn cancel failed: ", err)
			}
		}

		// IMPORTANT: 必须清空 Notify，否则死连接不会释放
		for err := range c.channelNotify {
			println(err)
		}
		for err := range c.connNotify {
			println(err)
		}
	}

quit:
	for {
		select {
		case <-c.quit:
			return
		default:
			log.Fatal("rabbitmq consumer - reconnect")

			if err := c.Run(); err != nil {
				log.Fatalf("rabbitmq consumer - failCheck:", err)

				// sleep 5s reconnect
				time.Sleep(time.Second * 5)
				continue
			}

			break quit
		}
	}
}
