package kafka

import (
	"context"

	"github.com/Shopify/sarama"
)

// Publish add data to queue
func Publish(ctx context.Context, name string, topic, msg string) error {
	p, err := DefaultManager.GetProducer(name)
	if err != nil {
		return err
	}
	return p.Publish(ctx, topic, msg)
}

// Consume data from queue
func Consume(ctx context.Context, name string, topics []string, handler sarama.ConsumerGroupHandler) error {
	c, err := DefaultManager.GetConsumer(name)
	if err != nil {
		return err
	}

	return c.Consume(ctx, topics, handler)
}

// ConsumeByPartition consume data by partition
func ConsumePartition(ctx context.Context, name, topic string, handler func([]byte) error) error {
	c, err := DefaultManager.GetConsumer(name)
	if err != nil {
		return err
	}

	return c.ConsumePartition(ctx, topic, handler)
}

// ConsumerByPartitionId consume data by partition id
func ConsumerByPartitionId(ctx context.Context, name, topic string, partition int32, handler func([]byte) error) error {
	c, err := DefaultManager.GetConsumer(name)
	if err != nil {
		return err
	}

	return c.ConsumerByPartitionId(ctx, topic, partition, handler)
}

// GetPartitionList get partition list
func GetPartitionList(ctx context.Context, name, topic string) ([]int32, error) {
	c, err := DefaultManager.GetConsumer(name)
	if err != nil {
		return nil, err
	}

	return c.client.Partitions(topic)
}
