package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

// OpenConnection connect to rabbitmq
func OpenConnection(addr string) (*amqp.Connection, error) {
	uri := fmt.Sprintf("amqp://%s", addr)

	return amqp.Dial(uri)
}
