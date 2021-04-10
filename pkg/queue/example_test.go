package queue

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/1024casts/snake/pkg/testing/lich"

	"github.com/Shopify/sarama"

	"github.com/1024casts/snake/pkg/queue/kafka"
	"github.com/1024casts/snake/pkg/queue/rabbitmq"
)

func TestMain(m *testing.M) {
	flag.Set("f", "../../test/rabbitmq-docker-compose.yaml")
	flag.Parse()

	if err := lich.Setup(); err != nil {
		panic(err)
	}
	defer lich.Teardown()

	if code := m.Run(); code != 0 {
		panic(code)
	}
}

func TestRabbitMQ(t *testing.T) {
	addr := "guest:guest@localhost:5672"
	conn, err := rabbitmq.OpenConnection(addr)
	if err != nil {
		t.Fatalf("failed connection: %s", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Fatalf("failed close connection: %s", err)
		}
	}()

	ch, err := rabbitmq.NewChannel(conn).Create()
	if err != nil {
		t.Fatalf("failed create channel: %s", err)
	}

	exchangeName := "test-exchange"
	queueName := "message-broker"

	err = ch.ExchangeDeclare(
		exchangeName,
		"direct",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("failed create exchange: %s", err)
	}

	var message = "Hello World RabbitMQ!"

	t.Run("rabbitmq publish message", func(t *testing.T) {
		producer := rabbitmq.NewProducer(addr, exchangeName)
		defer producer.Stop()
		if err := producer.Start(); err != nil {
			t.Errorf("start producer err: %s", err.Error())
		}
		if err := producer.Publish(message); err != nil {
			t.Errorf("failed publish message: %s", err.Error())
		}
	})

	// 自定义消息处理函数
	handler := func(body []byte) error {
		fmt.Println("consumer handler receive msg: ", string(body))
		return nil
	}

	t.Run("rabbitmq consume message", func(t *testing.T) {
		consumer := rabbitmq.NewConsumer(addr, exchangeName, queueName, true, handler)
		defer consumer.Stop()
		if err := consumer.Start(); err != nil {
			t.Errorf("failed consume: %s", err)
		}
	})
}

// TODO: read config
func TestKafka(t *testing.T) {
	var (
		config  = sarama.NewConfig()
		logger  = log.New(os.Stderr, "[sarama_logger]", log.LstdFlags)
		groupID = "sarama_consumer"
		topic   = "go-message-broker-topic"
		brokers = []string{"localhost:9093"}
		message = "Hello World Kafka!"
	)

	t.Run("kafka publish message", func(t *testing.T) {
		kafka.NewProducer(config, logger, topic, brokers).Publish(message)
	})

	t.Run("kafka consume message", func(t *testing.T) {
		kafka.NewConsumer(config, logger, topic, groupID, brokers).Consume()
	})
}
