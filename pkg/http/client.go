package http

import (
	"log"
)

// 禁止直接调用resty，统一使用HttpClient
var HttpClient = New("resty")

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
