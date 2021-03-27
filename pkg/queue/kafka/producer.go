package kafka

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"
)

type Producer struct {
	asyncProducer sarama.AsyncProducer
	topic         string
	enqueued      int
}

func NewProducer(config *sarama.Config, logger *log.Logger, topic string, brokers []string) *Producer {
	sarama.Logger = logger

	// Start a new async producer
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}

	log.Println("Kafka AsyncProducer up and running!")

	return &Producer{
		asyncProducer: producer,
		topic:         topic,
	}
}

func (p *Producer) Publish(message string) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		time.Sleep(5 * time.Second)
		message := &sarama.ProducerMessage{Topic: p.topic, Value: sarama.StringEncoder(message)}

		select {
		case p.asyncProducer.Input() <- message:
			p.enqueued++
			log.Printf("New message publish:  %s", message.Value)
		case <-signals:
			p.asyncProducer.AsyncClose() // Trigger a shutdown of the producer.
			log.Printf("Kafka AsyncProducer finished with %d messages produced.", p.enqueued)
			return
		}
	}

}
