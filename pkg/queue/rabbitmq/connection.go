package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

func OpenConnection(addr string) (*amqp.Connection, error) {
	uri := fmt.Sprintf("amqp://%s", addr)

	return amqp.Dial(uri)
}
