// http 客户端

package http

import (
	"context"
	"log"
	"time"
)

// see: https://github.com/iiinsomnia/gochat/blob/master/utils/http.go

const (
	contentTypeJSON = "application/json"
	contentTypeForm = "application/x-www-form-urlencoded"
)

// DefaultClient 默认的http client，基于resty库进行封装
var DefaultClient = "resty"

// RawClient 原生http client
var RawClient = "raw"

// Client 定义 http client 接口
type Client interface {
	Get(ctx context.Context, url string, params map[string]string, duration time.Duration, out interface{}) error
	Post(ctx context.Context, url string, data []byte, duration time.Duration, out interface{}) error
}

// New 实例化一个client
func New(opts ...Option) Client {
	cfg := option{
		ClientTyp: DefaultClient,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	var c Client
	if cfg.ClientTyp == DefaultClient {
		c = newRestyClient()
	} else {
		// c = newRawClient()
	}

	if c == nil {
		panic("unknown http client type " + cfg.ClientTyp)
	}

	log.Println(cfg.ClientTyp, "ready to serve")
	return c
}
