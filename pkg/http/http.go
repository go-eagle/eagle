package http

import "time"

type Client interface {
	Get(url string, params map[string]string, duration time.Duration) ([]byte, error)
	Post(url string, requestBody string, duration time.Duration) ([]byte, error)
}
