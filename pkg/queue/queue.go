package queue

import (
	"log"
	"os"

	"github.com/1024casts/snake/pkg/queue/rabbitmq"
)

const (
	PRODUCER = "producer"
	CONSUMER = "consumer"

	RABBITMQ = "rabbitmq"
	KAFKA    = "kafka"
	ACTIVEMQ = "activemq"
)

func NewRabbitMQ() {
	conn, err := rabbitmq.OpenConnection()
	if err != nil {
		log.Fatalf("failed connection: %s", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("failed close connection: %s", err)
		}
	}()

	channel, err := rabbitmq.NewChannel(conn).Create()
	if err != nil {
		log.Fatalf("failed create channel: %s", err)
	}

	queue, err := rabbitmq.NewQueue(channel, "go-message-broker").Create()
	if err != nil {
		log.Fatalf("failed queue declare: %s", err)
	}

	var message = "Hello World RabbitMQ!"

	switch os.Args[2] {
	case PRODUCER:
		if err := rabbitmq.NewProducer(channel, queue.Name).Publish(message); err != nil {
			log.Fatalf("failed publish message: %s", err)
		}
	case CONSUMER:
		if err := rabbitmq.NewConsumer(channel, queue.Name).Consume(); err != nil {
			log.Fatalf("failed consume: %s", err)
		}
	}
}
