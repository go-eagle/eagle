package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-eagle/eagle/pkg/queue/rabbitmq"
)

func main() {
	addr := "guest:guest@localhost:5672"

	// NOTE: need to create exchange and queue manually, than bind exchange to queue
	exchangeName := "test-exchange"
	queueName := "test-queue"

	// 自定义消息处理函数
	handler := func(ctx context.Context, body []byte) error {
		fmt.Println("consumer handler receive msg: ", string(body))
		return nil
	}

	// rabbitmq consume message
	ctx := context.Background()
	srv := rabbitmq.NewServer(addr, exchangeName)
	defer srv.Stop(ctx)

	err := srv.RegisterSubscriber(ctx, queueName, handler)
	if err != nil {
		log.Fatalf("RegisterSubscriber, err: %#v", err)
	}
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Start, err: %#v", err)
	}
}
