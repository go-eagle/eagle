// http 客户端

package http

import (
	"log"
	"time"
)

// see: https://github.com/iiinsomnia/gochat/blob/master/utils/http.go

const (
	contentTypeJSON = "application/json"
)

// DefaultClient 默认的http client，基于resty库进行封装
var DefaultClient = New("resty")

// RawClient 原生http client
var RawClient = New("raw")

// Client 定义 http client 接口
type Client interface {
	Get(url string, params map[string]string, duration time.Duration) ([]byte, error)
	Post(url string, data []byte, duration time.Duration) ([]byte, error)
}

// New 实例化一个client, default is raw http client
func New(typ string) Client {
	var c Client
	if typ == "resty" {
		c = newRestyClient()
	} else {
		c = newRawClient()
	}

	if c == nil {
		panic("unknown http client type " + typ)
	}

	log.Println(typ, "ready to serve")
	return c
}
