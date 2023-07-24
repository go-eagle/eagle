package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/config"
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
// go run examples/queue/rabbitmq/publisher.go -e local -c config
func main() {
	c := config.New(*cfgDir, config.WithEnv(*env))
	var cfg eagle.Config
	if err := c.Load("app", &cfg); err != nil {
		panic(err)
	}
	// set global
	eagle.Conf = &cfg

	rabbitmq.Load()
	defer rabbitmq.Close()

	var message string
	for i := 0; i < 10000; i++ {
		message = "Hello World RabbitMQ!" + time.Now().String()
		msg := map[string]interface{}{
			"message": message,
		}
		data, _ := json.Marshal(msg)
		if err := rabbitmq.Publish(context.Background(), "test", data, 0); err != nil {
			log.Fatalf("failed publish message: %s", err.Error())
		}
	}

}
