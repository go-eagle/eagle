package nats

import "log"

type Producer struct {
	topic string
}

func NewProducer(logger *log.Logger, topic string) *Producer {
	return &Producer{
		topic: topic,
	}
}

func (p *Producer) Publish(message string) {

}
