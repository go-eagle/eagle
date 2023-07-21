package options

import (
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// PublishOption defines the option for publishing
type PublishOption func(*PublishOptions)

// PublishOptions defines the options for publishing
type PublishOptions struct {
	amqp091.Publishing

	// Mandatory fails to publish if there are no queues bound to the routing key
	// if msg is not important, set mandatory to false, can improve system throughput
	Mandatory bool
	// Immediate fails to publish if there are no consumers  that can ack bound to the queue on the routing key
	// if msg is not important, set mandatory to false, can improve system throughput
	Immediate bool

	MsgExchange   string // if set, will override the exchange in producer
	MsgRoutingKey string // if set, will override the routing key in producer
	MaxRetry      uint64 // max retry times, default 0
}

// NewPublishOption returns a new PublishOptions
func NewPublishOption(opts ...PublishOption) *PublishOptions {
	options := PublishOptions{
		Publishing: amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
		},
	}
	for _, opt := range opts {
		opt(&options)
	}
	return &options
}

// WithPublishOptionMandatory sets the mandatory option
func WithPublishOptionMandatory(mandatory bool) PublishOption {
	return func(o *PublishOptions) {
		o.Mandatory = mandatory
	}
}

// WithPublishOptionImmediate sets the immediate option
func WithPublishOptionImmediate(immediate bool) PublishOption {
	return func(o *PublishOptions) {
		o.Immediate = immediate
	}
}

// WithPublishOptionMaxRetry sets the max retry option
func WithPublishOptionMaxRetry(MaxRetry uint64) PublishOption {
	return func(o *PublishOptions) {
		o.MaxRetry = MaxRetry
	}
}

// WithPublishOptionMsgExchange sets the msg exchange option
func WithPublishOptionMsgExchange(msgExchange string) PublishOption {
	return func(o *PublishOptions) {
		o.MsgExchange = msgExchange
	}
}

// WithPublishOptionMsgRoutingKey sets the msg routing key option
func WithPublishOptionMsgRoutingKey(msgRoutingKey string) PublishOption {
	return func(o *PublishOptions) {
		o.MsgRoutingKey = msgRoutingKey
	}
}

// WithPublishOptionHeaders sets the headers option
func WithPublishOptionHeaders(headers amqp091.Table) PublishOption {
	return func(o *PublishOptions) {
		o.Headers = headers
	}
}

// WithPublishOptionContentType sets the content type option
func WithPublishOptionContentType(contentType string) PublishOption {
	return func(o *PublishOptions) {
		o.ContentType = contentType
	}
}

// WithPublishOptionContentEncoding sets the content type option
func WithPublishOptionContentEncoding(contentEncoding string) PublishOption {
	return func(o *PublishOptions) {
		o.ContentEncoding = contentEncoding
	}
}

// WithPublishOptionDeliveryMode sets the delivery mode option
func WithPublishOptionDeliveryMode(deliveryMode uint8) PublishOption {
	return func(o *PublishOptions) {
		o.DeliveryMode = deliveryMode
	}
}

// WithPublishOptionPriority sets the priority option
func WithPublishOptionPriority(priority uint8) PublishOption {
	return func(o *PublishOptions) {
		o.Priority = priority
	}
}

// WithPublishOptionCorrelationID sets the correlation id option
func WithPublishOptionCorrelationID(correlationID string) PublishOption {
	return func(o *PublishOptions) {
		o.CorrelationId = correlationID
	}
}

// WithPublishOptionReplyTo sets the reply to option
func WithPublishOptionReplyTo(replyTo string) PublishOption {
	return func(o *PublishOptions) {
		o.ReplyTo = replyTo
	}
}

// WithPublishOptionExpiration sets the expiration option
func WithPublishOptionExpiration(expiration string) PublishOption {
	return func(o *PublishOptions) {
		o.Expiration = expiration
	}
}

// WithPublishOptionMessageID sets the message id option
func WithPublishOptionMessageID(messageID string) PublishOption {
	return func(o *PublishOptions) {
		o.MessageId = messageID
	}
}

// WithPublishOptionTimestamp sets the timestamp option
func WithPublishOptionTimestamp(timestamp time.Time) PublishOption {
	return func(o *PublishOptions) {
		o.Timestamp = timestamp
	}
}

// WithPublishOptionType sets the type option
func WithPublishOptionType(typ string) PublishOption {
	return func(o *PublishOptions) {
		o.Type = typ
	}
}

// WithPublishOptionUserID sets the user id option
func WithPublishOptionUserID(userID string) PublishOption {
	return func(o *PublishOptions) {
		o.UserId = userID
	}
}

// WithPublishOptionAppID sets the app id option
func WithPublishOptionAppID(appID string) PublishOption {
	return func(o *PublishOptions) {
		o.AppId = appID
	}
}
