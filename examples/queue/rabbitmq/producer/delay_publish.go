package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/config"
	logger "github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/queue/rabbitmq"
	"github.com/go-eagle/eagle/pkg/queue/rabbitmq/options"
)

// 启动 rabbitmq
// docker run -it  --name rabbitmq -p 5672:5672 -p 15672:15672 -v $PWD/plugins:/plugins rabbitmq:3.10-management
// 访问ui: http://127.0.0.1:15672/
// cd examples/queue/rabbitmq/producer
// go run delay_publish.go
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

	var message string
	for i := 0; i < 100000; i++ {
		message = "Hello World RabbitMQ!" + time.Now().String()
		msg := map[string]interface{}{
			"message": message,
		}
		data, _ := json.Marshal(msg)
		if err := rabbitmq.PublishWithDelay(context.Background(), "test-demo", data, 10, opts...); err != nil {
			log.Fatalf("failed publish message: %s", err.Error())
		}
	}
}
