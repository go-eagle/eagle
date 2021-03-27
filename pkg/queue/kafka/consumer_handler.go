package kafka

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

// ConsumerGroupHandler represents the sarama consumer group
type ConsumerGroupHandler struct{}

// Setup is run before consumer start consuming, is normally used to setup things such as database connections
func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages(), here is supposed to be what you want to
// do with the message. In this example the message will be logged with the topic name, partition and message value.
func (h ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf(
			"Message topic:%q partition:%d offset:%d message: %v\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value),
		)

		session.MarkMessage(msg, "")

		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	}

	return nil
}
