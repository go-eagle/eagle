package rabbitmq

import "github.com/streadway/amqp"

type Queue struct {
	channel *amqp.Channel
	name    string
}

func NewQueue(channel *amqp.Channel, name string) Queue {
	return Queue{channel: channel, name: name}
}

func (q Queue) Create() (amqp.Queue, error) {
	return q.channel.QueueDeclare(
		q.name,
		false,
		false,
		false,
		false,
		nil,
	)
}
