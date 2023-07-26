package main

import (
	"context"
	"encoding/json"

	"github.com/go-eagle/eagle/pkg/config"
	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/pflag"

	eagle "github.com/go-eagle/eagle/pkg/app"
	logger "github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/queue/rabbitmq"
	"github.com/go-eagle/eagle/pkg/queue/rabbitmq/options"
	RabbitMQ "github.com/go-eagle/eagle/pkg/transport/consumer/rabbitmq"
)

// cd examples/queue/rabbitmq/consumer/
// go run server.go
func main() {
	pflag.Parse()

	// init config
	c := config.New("config")
	var cfg eagle.Config
	if err := c.Load("app", &cfg); err != nil {
		panic(err)
	}
	// set global
	eagle.Conf = &cfg

	logger.Init()

	rabbitmq.Load()
	defer rabbitmq.Close()

	// 自定义消息处理函数
	handler := func(ctx context.Context, body amqp091.Delivery) (action rabbitmq.Action) {
		msg := make(map[string]interface{})
		err := json.Unmarshal(body.Body, &msg)
		if err != nil {
			logger.Errorf("consumer handler unmarshal msg err: %s", err.Error())
			return rabbitmq.NackDiscard
		}
		logger.Infof("consumer handler receive msg: %s", msg)
		return rabbitmq.Ack
	}
	handler2 := func(ctx context.Context, body amqp091.Delivery) (action rabbitmq.Action) {
		msg := make(map[string]interface{})
		err := json.Unmarshal(body.Body, &msg)
		if err != nil {
			logger.Errorf("consumer handler unmarshal msg err: %s", err.Error())
			return rabbitmq.NackDiscard
		}
		logger.Infof("consumer handler receive msg: %s", msg)
		return rabbitmq.Ack
	}

	// rabbitmq consume message
	opts := []options.ConsumerOption{
		options.WithConsumerOptionConcurrency(1),
	}

	srv := RabbitMQ.NewServer(opts...)

	// register subscriber can place into init function in internal/task/task.go
	err := srv.RegisterHandler("test-demo", handler)
	if err != nil {
		panic(err)
	}
	err = srv.RegisterHandler("test-multi", handler2)
	if err != nil {
		panic(err)
	}

	// start app
	app := eagle.New(
		eagle.WithName(cfg.Name),
		eagle.WithVersion(cfg.Version),
		eagle.WithLogger(logger.GetLogger()),
		eagle.WithServer(
			srv,
		),
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
