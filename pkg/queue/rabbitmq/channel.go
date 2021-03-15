package rabbitmq

import "github.com/streadway/amqp"

type Channel struct {
	conn *amqp.Connection
}

func NewChannel(conn *amqp.Connection) Channel {
	return Channel{conn: conn}
}

func (c Channel) Create() (*amqp.Channel, error) {
	return c.conn.Channel()
}
