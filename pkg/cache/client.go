package cache

import "log"

// 禁止直接调用redis，统一使用该变量
var CacheClient = New("redis")

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
