package main

import (
	"log"
	"time"

	"github.com/go-eagle/eagle/pkg/queue/rabbitmq"
)

func main() {
	addr := "guest:guest@localhost:5672"

	// NOTE: need to create exchange and queue manually, than bind exchange to queue
	exchangeName := "test-exchange"
	// like topic, bind to queue: test-queue
	routingKey := "test-routing-key"

	// rabbitmq publish message
	producer := rabbitmq.NewProducer(addr, exchangeName)
	defer producer.Stop()
	if err := producer.Start(); err != nil {
		log.Fatalf("start producer err: %s", err.Error())
	}

	var message string
	for i := 0; i < 10000; i++ {
		message = "Hello World RabbitMQ!" + time.Now().String()
		if err := producer.Publish(routingKey, message); err != nil {
			log.Fatalf("failed publish message: %s", err.Error())
		}
	}

}
