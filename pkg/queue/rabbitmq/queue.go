package rabbitmq

import "github.com/streadway/amqp"

// Queue .
type Queue struct {
	channel *amqp.Channel
	name    string
}

// NewQueue .
func NewQueue(channel *amqp.Channel, name string) Queue {
	return Queue{channel: channel, name: name}
}

// Create .
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
