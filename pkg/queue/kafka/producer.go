package kafka

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"

	logger "github.com/go-eagle/eagle/pkg/log"
)

// Producer kafka producer
type Producer struct {
	asyncProducer sarama.AsyncProducer
	topic         string
	enqueued      int
}

// NewProducer create producer
// nolint
func NewProducer(config *sarama.Config, logger *log.Logger, topic string, brokers []string) *Producer {
	sarama.Logger = logger

	// Start a new async producer
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}

	log.Println("Kafka AsyncProducer up and running!")

	p := &Producer{
		asyncProducer: producer,
		topic:         topic,
	}

	go p.asyncDealMessage()

	return p
}

func (p *Producer) asyncDealMessage() {
	for {
		select {
		case res := <-p.asyncProducer.Successes():
			logger.Info("push msg success", "topic is", res.Topic, "partition is ", res.Partition, "offset is ", res.Offset)
		case err := <-p.asyncProducer.Errors():
			logger.Info("push msg failed", "err is ", err.Error())
		}
	}
}

// Publish push data to queue
func (p *Producer) Publish(message string) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		time.Sleep(5 * time.Second)
		message := &sarama.ProducerMessage{Topic: p.topic, Value: sarama.StringEncoder(message)}

		select {
		case p.asyncProducer.Input() <- message:
			p.enqueued++
			logger.Infof("New message publish:  %s", message.Value)
		case <-signals:
			p.asyncProducer.AsyncClose() // Trigger a shutdown of the producer.
			logger.Infof("Kafka AsyncProducer finished with %d messages produced.", p.enqueued)
			return
		}
	}
}
