package http

import (
	"testing"
	"time"
)

func TestServer(t *testing.T) {

	srv := NewServer()

	//go func() {
	if err := srv.Start(); err != nil {
		panic(err)
	}
	//}()
	time.Sleep(time.Second)
	srv.Stop()
}
