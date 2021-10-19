package main

import (
	"encoding/base64"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/go-eagle/eagle/pkg/queue/nats"
)

func main() {
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
					log.Fatal(err)
				}
				log.Println("producer handler publish msg: ", "hello nats")

			case <-received:
				wg.Done()
				break
			}
		}
	}()
	go func() {
		for {
			// nolint: gosimple
			select {
			default:
				handler := func(message []byte) error {
					decodeMessage, _ := base64.StdEncoding.DecodeString(strings.Trim(string(message), "\""))
					log.Println("consumer handler receive msg: ", string(decodeMessage))
					received <- struct{}{}
					wg.Done()
					return nil
				}
				if err := consumer.Consume(topic, handler); err != nil {
					log.Fatal(err)
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()

	published <- struct{}{}
	wg.Wait()
}
