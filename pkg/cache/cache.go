package cache

import (
	"log"
	"time"
)

type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (interface{}, error)
	Del(key string) error
}

// 禁止直接调用redis，统一使用该变量
var RedisCache = New("redis")

func New(typ string) Cache {
	var c Cache
	if typ == "redis" {
		c = newRedisCache()
	}

	if c == nil {
		panic("unknown cache type " + typ)
	}

	log.Println(typ, "ready to serve")
	return c
}
