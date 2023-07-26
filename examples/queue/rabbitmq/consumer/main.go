package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-eagle/eagle/pkg/queue/rabbitmq/options"

	"github.com/rabbitmq/amqp091-go"

	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/config"

	"github.com/spf13/pflag"

	logger "github.com/go-eagle/eagle/pkg/log"

	"github.com/go-eagle/eagle/pkg/queue/rabbitmq"
)

var (
	cfgDir = pflag.StringP("config dir", "c", "config", "config path.")
	env    = pflag.StringP("env name", "e", "", "env var name.")
)

// cd examples/queue/rabbitmq/consumer/
// go run main.go
func main() {
	pflag.Parse()

	// init config
	c := config.New(*cfgDir, config.WithEnv(*env))
	var cfg eagle.Config
	if err := c.Load("app", &cfg); err != nil {
		panic(err)
	}
	// set global
	eagle.Conf = &cfg

	logger.Init()

	rabbitmq.Load()
	defer rabbitmq.Close()

	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})
	stop := make(chan struct{}, 1)

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

	// rabbitmq consume message
	ctx := context.Background()

	opts := []options.ConsumerOption{
		options.WithConsumerOptionConcurrency(1),
	}

	go func() {
		err := rabbitmq.Consume(ctx, "test-demo", handler, opts...)
		if err != nil {
			logger.Errorf("rabbitmq consume err: %s", err.Error())
		}
	}()

	for {
		select {
		case <-stopSig:
			logger.Info("received stop signal")
			stop <- struct{}{}
		case <-stop:
			logger.Info("stopping service")
			close(done)
			return
		case <-done:
			logger.Info("stopped service gracefully")
			return
		}
	}
}
