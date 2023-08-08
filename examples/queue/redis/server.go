package main

import (
	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/config"
	logger "github.com/go-eagle/eagle/pkg/log"
	redisMQ "github.com/go-eagle/eagle/pkg/transport/consumer/redis"
	"github.com/hibiken/asynq"
	"github.com/spf13/pflag"
)

// redis queue consumer
// cd examples/queue/redis/consumer/
// go run server.go handler.go client.go
func main() {
	pflag.Parse()

	// init config
	c := config.New(".")
	var cfg eagle.Config
	if err := c.Load("app", &cfg); err != nil {
		panic(err)
	}
	// set global
	eagle.Conf = &cfg

	logger.Init()

	srv := redisMQ.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				redisMQ.QueueCritical: 6,
				redisMQ.QueueDefault:  3,
				redisMQ.QueueLow:      1,
			},
			// See the godoc for other configuration options
		},
	)

	// register handler
	srv.RegisterHandler(TypeEmailWelcome, HandleEmailWelcomeTask)
	// here register other handlers...

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
