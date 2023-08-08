package main

import "time"

// cd examples/queue/redis/consumer/
// go run producer.go handler.go client.go
func main() {
	for i := 0; i < 10; i++ {
		err := NewEmailWelcomeTask(EmailWelcomePayload{
			UserID: time.Now().Unix(),
		})
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
}
