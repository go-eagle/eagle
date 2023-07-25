package rabbitmq

import (
	"context"
	"fmt"

	"github.com/go-eagle/eagle/pkg/queue/rabbitmq/options"
)

// Publish add data to queue
func Publish(ctx context.Context, name string, data []byte, retry uint64, opts ...options.PublishOption) error {
	p, err := DefaultManager.GetProducer(name)
	if err != nil {
		return err
	}
	return p.Publish(ctx, data, opts...)
}

// Consume data from queue
func Consume(ctx context.Context, name string, handler Handler, opts ...options.ConsumerOption) error {
	c, err := DefaultManager.GetConsumer(name)
	if err != nil {
		return err
	}
	if c.IsClosed() {
		return fmt.Errorf("rabbitmq: consumer %s closed", name)
	}
	return c.Consume(ctx, handler, opts...)
}
