package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	logger "github.com/go-eagle/eagle/pkg/log"
)

// Producer kafka producer
type Producer struct {
	asyncProducer sarama.AsyncProducer
	enqueued      int
	logger        logger.Logger
}

// NewProducer create producer
// nolint
func NewProducer(conf *Conf, logger logger.Logger) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.RequiredAcks(conf.RequiredAcks)
	config.Producer.Retry.Max = 3
	config.Producer.Partitioner = getPartitoner(conf)

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("kafka: producer config validate error: %v", err)
	}

	// Start a new async producer
	producer, err := sarama.NewAsyncProducer(conf.Brokers, config)
	if err != nil {
		return nil, err
	}

	log.Println("kafka: AsyncProducer up and running!")

	p := &Producer{
		asyncProducer: producer,
		logger:        logger,
	}

	go p.asyncDealMessage()

	return p, nil
}

func (p *Producer) asyncDealMessage() {
	for {
		select {
		case res := <-p.asyncProducer.Successes():
			p.logger.Info("kafka: push msg success", "topic is", res.Topic, "partition is ", res.Partition, "offset is ", res.Offset)
		case err := <-p.asyncProducer.Errors():
			p.logger.Error("kafka: push msg failed", "err is ", err.Error())
		}
	}
}

// Publish push data to queue
func (p *Producer) Publish(ctx context.Context, topic string, message string) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		message := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(message),
		}

		select {
		case p.asyncProducer.Input() <- message:
			p.enqueued++
			p.logger.Infof("kafka: New message publish:  %s", message.Value)
		case <-signals:
			p.asyncProducer.AsyncClose() // Trigger a shutdown of the producer.
			p.logger.Infof("kafka: AsyncProducer finished with %d messages produced.", p.enqueued)
			return nil
		}
	}
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.asyncProducer.Close()
}

// getPartitoner returns the partitioner constructor based on the configuration
func getPartitoner(conf *Conf) sarama.PartitionerConstructor {
	switch conf.Partitioner {
	case "random":
		return sarama.NewRandomPartitioner
	case "roundrobin":
		return sarama.NewRoundRobinPartitioner
	case "hash":
		return sarama.NewHashPartitioner
	default:
		return sarama.NewRandomPartitioner
	}
}
