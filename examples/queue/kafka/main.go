package main

import (
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/go-eagle/eagle/pkg/queue/kafka"
)

func main() {
	var (
		config  = sarama.NewConfig()
		logger  = log.New(os.Stderr, "[sarama_logger]", log.LstdFlags)
		groupID = "sarama_consumer"
		topic   = "go-message-broker-topic"
		brokers = []string{"localhost:9093"}
		message = "Hello World Kafka!"
	)

	// kafka publish message
	kafka.NewProducer(config, logger, topic, brokers).Publish(message)

	// kafka consume message
	kafka.NewConsumer(config, logger, topic, groupID, brokers).Consume()

}
