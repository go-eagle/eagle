package http

import "time"

// http client 接口
type Client interface {
	Get(url string, params map[string]string, duration time.Duration) ([]byte, error)
	Post(url string, requestBody string, duration time.Duration) ([]byte, error)
}
