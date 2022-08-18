package main

import (
	"context"

	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/config"

	"github.com/spf13/pflag"

	logger "github.com/go-eagle/eagle/pkg/log"

	"github.com/go-eagle/eagle/pkg/queue/rabbitmq"
)

var (
	cfgDir  = pflag.StringP("config dir", "c", "config", "config path.")
	env     = pflag.StringP("env name", "e", "", "env var name.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

// go run examples/queue/rabbitmq/consumer.go -e local -c config
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

	addr := "guest:guest@localhost:5672"

	// NOTE: need to create exchange and queue manually, than bind exchange to queue
	exchangeName := "test-exchange"
	queueName := "test-queue"
	queueName2 := "test-queue2"

	done := make(chan struct{})

	// 自定义消息处理函数
	handler := func(ctx context.Context, body []byte) error {
		logger.Infof("consumer handler receive msg: %s", string(body))
		return nil
	}
	handler2 := func(ctx context.Context, body []byte) error {
		logger.Infof("consumer handler2 receive msg: %s", string(body))
		return nil
	}

	// rabbitmq consume message
	ctx := context.Background()
	srv := rabbitmq.NewServer(addr, exchangeName)
	defer srv.Stop(ctx)

	err := srv.RegisterSubscriber(ctx, queueName, handler)
	if err != nil {
		logger.Errorf("RegisterSubscriber, err: %#v", err)
	}
	err = srv.RegisterSubscriber(ctx, queueName2, handler2)
	if err != nil {
		logger.Errorf("RegisterSubscriber, err: %#v", err)
	}

	if err := srv.Start(ctx); err != nil {
		logger.Errorf("Start, err: %#v", err)
	}

	<-done

}
