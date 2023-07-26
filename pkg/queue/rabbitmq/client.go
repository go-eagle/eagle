package rabbitmq

import (
	"context"
	"fmt"

	"github.com/go-eagle/eagle/pkg/queue/rabbitmq/options"
)

// Publish add data to queue
func Publish(ctx context.Context, name string, data []byte, opts ...options.PublishOption) error {
	p, err := DefaultManager.GetProducer(name)
	if err != nil {
		return err
	}
	return p.Publish(ctx, data, opts...)
}

// PublishWithDelay add a delay msg to queue
// delayTime: seconds
func PublishWithDelay(ctx context.Context, name string, data []byte, delayTime int, opts ...options.PublishOption) error {
	p, err := DefaultManager.GetProducer(name)
	if err != nil {
		return err
	}
	opts = append(opts, options.WithPublishOptionHeaders(map[string]interface{}{
		"x-delay": delayTime * 1000, // seconds
	}))
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
