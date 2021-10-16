package http

import (
	"context"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	srv := NewServer()

	//go func() {
	if err := srv.Start(context.Background()); err != nil {
		panic(err)
	}
	//}()
	time.Sleep(time.Second)
	srv.Stop(context.Background())
}
