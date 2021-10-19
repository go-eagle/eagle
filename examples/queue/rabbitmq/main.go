package main

import (
	"fmt"
	"log"

	"github.com/go-eagle/eagle/pkg/queue/rabbitmq"
)

func main() {
	addr := "guest:guest@localhost:5672"

	// NOTE: need to create exchange and queue manually, than bind exchange to queue
	exchangeName := "test-exchange"
	queueName := "test-bind-to-exchange"

	var message = "Hello World RabbitMQ!"

	// rabbitmq publish message
	producer := rabbitmq.NewProducer(addr, exchangeName)
	defer producer.Stop()
	if err := producer.Start(); err != nil {
		log.Fatalf("start producer err: %s", err.Error())
	}
	if err := producer.Publish(message); err != nil {
		log.Fatalf("failed publish message: %s", err.Error())
	}

	// 自定义消息处理函数
	handler := func(body []byte) error {
		fmt.Println("consumer handler receive msg: ", string(body))
		return nil
	}

	// rabbitmq consume message
	// NOTE: autoDelete param
	consumer := rabbitmq.NewConsumer(addr, exchangeName, queueName, false, handler)
	defer consumer.Stop()
	if err := consumer.Start(); err != nil {
		log.Fatalf("failed consume: %s", err)
	}
}
