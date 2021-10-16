package rabbitmq

import "github.com/streadway/amqp"

// Channel data channel
type Channel struct {
	conn *amqp.Connection
}

// NewChannel instance a channel
func NewChannel(conn *amqp.Connection) Channel {
	return Channel{conn: conn}
}

// Create create a channel
func (c Channel) Create() (*amqp.Channel, error) {
	return c.conn.Channel()
}
