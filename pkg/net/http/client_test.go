package http

import (
	"testing"
	"time"
)

func TestDefaultClient(t *testing.T) {
	c := DefaultClient

	t.Run("test http get func", func(t *testing.T) {
		b, err := c.Get("http://httpbin.org/get", nil, 3*time.Second)
		if err != nil {
			t.Error(err)
		}

		t.Log(string(b))
	})

	t.Run("test http post func", func(t *testing.T) {
		b, err := c.Post("http://httpbin.org/post", nil, 3*time.Second)
		if err != nil {
			t.Error(err)
		}

		t.Log(string(b))
	})
}
