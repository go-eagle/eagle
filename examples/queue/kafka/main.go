package main

import (
	"context"
	"log"

	"github.com/go-eagle/eagle/pkg/queue/kafka"
)

func main() {
	// 1. 初始化配置
	kafka.Load()
	defer kafka.Close()

	// 2. 获取配置信息（可选）
	configs := kafka.GetConfig()
	if len(configs) == 0 {
		log.Fatal("No kafka config found")
	}

	// 3. 使用配置进行消息发布
	ctx := context.Background()
	err := kafka.Publish(ctx, "default", "test-topic", "hello world")
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}

	// 4. 使用配置进行消息消费
	handler := func(data []byte) error {
		log.Printf("Received message: %s", string(data))
		return nil
	}

	// 从默认实例消费
	go func() {
		err := kafka.ConsumePartition(ctx, "default", "test-topic", handler)
		if err != nil {
			log.Printf("Failed to consume message: %v", err)
		}
	}()

	// 从order实例消费
	go func() {
		err := kafka.ConsumePartition(ctx, "order", "order-topic", handler)
		if err != nil {
			log.Printf("Failed to consume message: %v", err)
		}
	}()

	select {}
}
