package queue

// Producer queue producer
type Producer interface {
	Publish(message string) error
}

// Consumer queue consumer
type Consumer interface {
	Consume() error
}
