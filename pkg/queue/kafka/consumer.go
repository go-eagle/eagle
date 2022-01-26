package kafka

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"

	logger "github.com/go-eagle/eagle/pkg/log"
)

// Consumer kafka consumer
type Consumer struct {
	group   sarama.ConsumerGroup
	topics  []string
	groupID string
	handler sarama.ConsumerGroupHandler

	ctx    context.Context
	cancel context.CancelFunc
}

// NewConsumer create a consumer
// nolint
func NewConsumer(config *sarama.Config, logger *log.Logger, topic string, groupID string, brokers []string, handler *ConsumerGroupHandler) *Consumer {
	// Init config, specify appropriate versio
	sarama.Logger = log.New(os.Stderr, "[sarama_logger]", log.LstdFlags)
	sarama.Logger = logger
	config.Version = sarama.V2_0_0_0 // V2_4_0_0

	// Start with a client
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		panic(err)
	}

	// Start a new consumer group
	group, err := sarama.NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Consumer{
		group:   group,
		topics:  []string{topic},
		groupID: groupID,
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Consume consume data
func (c *Consumer) Consume() {
	// Track errors
	go func() {
		for err := range c.group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		select {
		case <-c.ctx.Done():
			_ = c.group.Close()
			logger.Info("[Kafka] Consume ctx done")
			return
		default:
			if err := c.group.Consume(ctx, c.topics, c.handler); err != nil {
				logger.Errorf("[Kafka] Consume err: %s", err.Error())
			}
		}
	}
}

// Stop close conn
func (c Consumer) Stop() {
	c.cancel()
}
