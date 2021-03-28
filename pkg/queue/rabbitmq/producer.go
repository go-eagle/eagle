package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Producer struct {
	channel   *amqp.Channel
	queueName string
}

func NewProducer(channel *amqp.Channel, queueName string) Producer {
	return Producer{channel: channel, queueName: queueName}
}

func (p Producer) Publish(message string) error {
	if err := p.channel.Publish(
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			Headers:     amqp.Table{},
			ContentType: "text/plain",
			Body:        []byte(message),
		}); err != nil {
		return fmt.Errorf("failed to publish a message: %s", err)
	}

	return nil
}
