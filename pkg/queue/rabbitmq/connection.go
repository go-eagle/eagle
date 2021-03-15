package rabbitmq

import "github.com/streadway/amqp"

func OpenConnection() (*amqp.Connection, error) {
	uri := "amqp://guest:guest@localhost:5672"

	return amqp.Dial(uri)
}
