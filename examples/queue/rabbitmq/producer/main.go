package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-eagle/eagle/pkg/queue/rabbitmq/options"

	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/config"
	logger "github.com/go-eagle/eagle/pkg/log"
	"github.com/spf13/pflag"

	"github.com/go-eagle/eagle/pkg/queue/rabbitmq"
)

var (
	cfgDir = pflag.StringP("config dir", "c", "config", "config path.")
	env    = pflag.StringP("env name", "e", "", "env var name.")
)

// 启动 rabbitmq
// docker run -it  --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.10-management
// 访问ui: http://127.0.0.1:15672/
// cd examples/queue/rabbitmq/producer
// go run main.go
func main() {
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

	opts := []options.PublishOption{
		options.WithPublishOptionContentType("application/json"),
	}

	go func() {
		var message string
		for i := 0; i < 100000; i++ {
			message = "Hello World RabbitMQ!" + time.Now().String()
			msg := map[string]interface{}{
				"message": message,
			}
			data, _ := json.Marshal(msg)
			if err := rabbitmq.Publish(context.Background(), "test-demo", data, opts...); err != nil {
				log.Fatalf("failed publish message: %s", err.Error())
			}
		}
	}()

	var message string
	for i := 0; i < 100000; i++ {
		message = "Hello World multi RabbitMQ!" + time.Now().String()
		msg := map[string]interface{}{
			"message": message,
		}
		data, _ := json.Marshal(msg)
		if err := rabbitmq.Publish(context.Background(), "test-multi", data, opts...); err != nil {
			log.Fatalf("failed publish message: %s", err.Error())
		}
	}

}
