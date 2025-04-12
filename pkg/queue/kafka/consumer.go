package kafka

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/go-eagle/eagle/pkg/log"
)

// Consumer kafka consumer
type Consumer struct {
	client   sarama.Client
	logger   log.Logger
	group    sarama.ConsumerGroup
	consumer sarama.Consumer
}

// NewConsumer create a consumer
// nolint
func NewConsumer(conf *Conf, logger log.Logger) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Retry.Backoff = 100
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	version, err := sarama.ParseKafkaVersion(conf.Version)
	if err != nil {
		return nil, fmt.Errorf("kafka: parse version %s error: %v", conf.Version, err)
	}
	config.Version = version

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("kafka: consumer config validate error: %v", err)
	}

	// Start with a client
	client, err := sarama.NewClient(conf.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("kafka: create client error: %v", err)
	}

	// Start a new consumer group
	group, err := sarama.NewConsumerGroupFromClient(conf.GroupID, client)
	if err != nil {
		return nil, fmt.Errorf("kafka: create consumer group error: %v", err)
	}

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, fmt.Errorf("kafka: create consumer from client error: %v", err)
	}

	return &Consumer{
		client:   client,
		logger:   logger,
		group:    group,
		consumer: consumer,
	}, nil
}

// Consume consume data
// This method consumes messages from the specified topics using a consumer group.
// It tracks errors in a separate goroutine and logs them.
// The handler function is called for each message received.
func (c *Consumer) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	// Track errors
	go func() {
		for err := range c.group.Errors() {
			c.logger.Errorf("kafka: Consume group errors: %v", err)
		}
	}()

	// Iterate over consumer sessions.
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("kafka: Consume ctx done")
			return nil
		default:
			if err := c.group.Consume(ctx, topics, handler); err != nil {
				c.logger.Errorf("kafka: Group Consume err: %v", err)
			}
		}
	}
}

// ConsumePartition consume data by partition
// This method consumes messages from all partitions of the specified topic.
// It creates a separate goroutine for each partition to handle messages concurrently.
// The handler function is called for each message received.
func (c *Consumer) ConsumePartition(ctx context.Context, topic string, handler func([]byte) error) error {
	partitionList, err := c.consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for _, partition := range partitionList {
		pc, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				if err := handler(msg.Value); err != nil {
					// handle error
					c.logger.Errorf("kafka: ConsumePartition message error, topic %s, partition %d, err: %v", topic, msg.Partition, err)
					continue
				}
			}
		}(pc)
	}

	return nil
}

// Consumer by partition id
func (c *Consumer) ConsumerByPartitionId(ctx context.Context, topic string, partition int32, handler func([]byte) error) error {
	pc, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		return nil
	}

	for msg := range pc.Messages() {
		if err := handler(msg.Value); err != nil {
			// handle error
			c.logger.Errorf("kafka: ConsumerByPartitionId message error, topic %s partition %d err: %v", topic, partition, err)
			continue
		}
	}
	return nil
}

// Stop close conn
func (c *Consumer) Stop() {
	if c.consumer != nil {
		c.consumer.Close()
	}
	if c.group != nil {
		c.group.Close()
	}
	if c.client != nil {
		c.client.Close()
	}
}
