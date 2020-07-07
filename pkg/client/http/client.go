// http 客户端

package http

import (
	"log"
)

// see: https://github.com/iiinsomnia/gochat/blob/master/utils/http.go

// HTTPClient 禁止直接调用resty，统一使用HttpClient
var HTTPClient = New("resty")

// New 实例化一个client
func New(typ string) Client {
	var c Client
	if typ == "resty" {
		c = newRestyClient()
	}

	if c == nil {
		panic("unknown http type " + typ)
	}

	log.Println(typ, "ready to serve")
	return c
}
