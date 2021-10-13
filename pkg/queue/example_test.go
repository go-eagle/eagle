package queue

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Shopify/sarama"

	"github.com/go-eagle/eagle/pkg/queue/kafka"
	"github.com/go-eagle/eagle/pkg/queue/nats"
	"github.com/go-eagle/eagle/pkg/queue/rabbitmq"
	"github.com/go-eagle/eagle/pkg/testing/lich"
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

	// NOTE: need to create exchange and queue manually, than bind exchange to queue
	exchangeName := "test-exchange"
	queueName := "test-bind-to-exchange"

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
		// NOTE: autoDelete param
		consumer := rabbitmq.NewConsumer(addr, exchangeName, queueName, false, handler)
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

func TestNats(t *testing.T) {
	var (
		addr  = "nats://localhost:4222"
		topic = "hello"
	)
	producer := nats.NewProducer(addr)
	consumer := nats.NewConsumer(addr)

	published := make(chan struct{})
	received := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for {
			select {
			case <-published:
				time.Sleep(3 * time.Second)
				if err := producer.Publish(topic, []byte("hello nats")); err != nil {
					t.Error(err)
				}
				t.Log("producer handler publish msg: ", "hello nats")

			case <-received:
				wg.Done()
				break
			}
		}
	}()
	go func() {
		for {
			select {
			default:
				handler := func(message []byte) error {
					decodeMessage, _ := base64.StdEncoding.DecodeString(strings.Trim(string(message), "\""))
					t.Log("consumer handler receive msg: ", string(decodeMessage))
					received <- struct{}{}
					wg.Done()
					return nil
				}
				if err := consumer.Consume(topic, handler); err != nil {
					t.Error(err)
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()

	published <- struct{}{}
	wg.Wait()
}
