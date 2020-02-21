package http

import "time"

type Client interface {
	Get(url string, duration time.Duration) ([]byte, error)
	Post(url string, requestBody string, duration time.Duration) ([]byte, error)
}
