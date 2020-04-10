package http

import "time"

// Client 定义 http client 接口
type Client interface {
	Get(url string, params map[string]string, duration time.Duration) ([]byte, error)
	Post(url string, requestBody string, duration time.Duration) ([]byte, error)
}
