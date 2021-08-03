package http

import (
	"context"
	"testing"
	"time"
)

func TestDefaultClient(t *testing.T) {
	c := DefaultClient

	t.Run("test http get func", func(t *testing.T) {
		var ret interface{}
		err := c.Get(context.Background(), "http://httpbin.org/get", nil, 3*time.Second, &ret)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(ret)
	})

	t.Run("test http post func", func(t *testing.T) {
		var ret interface{}
		err := c.Post(context.Background(), "http://httpbin.org/post", nil, 3*time.Second, &ret)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(ret)
	})
}

func TestRawClient(t *testing.T) {
	c := RawClient

	t.Run("test http get func", func(t *testing.T) {
		var ret interface{}
		err := c.Get(context.Background(), "http://httpbin.org/get", nil, 3*time.Second, &ret)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(ret)
	})

	t.Run("test http post func", func(t *testing.T) {
		var ret interface{}
		err := c.Post(context.Background(), "http://httpbin.org/post", nil, 3*time.Second, &ret)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(ret)
	})
}
