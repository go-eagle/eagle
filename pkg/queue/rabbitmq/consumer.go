package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	channel   *amqp.Channel
	queueName string
}

func NewConsumer(channel *amqp.Channel, queueName string) Consumer {
	return Consumer{channel: channel, queueName: queueName}
}

func (c Consumer) Consume() error {
	deliveries, err := c.channel.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue consume: %s", err)
	}

	done := make(chan bool)

	go func() {
		for d := range deliveries {
			log.Printf("Consumer received a message: %s in queue: %s", d.Body, c.queueName)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-done

	return nil
}

func handler(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		d.Ack(false)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
