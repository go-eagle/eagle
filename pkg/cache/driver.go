package cache

import (
	"time"

	"github.com/gin-gonic/gin"

	redis2 "github.com/1024casts/snake/pkg/redis"
)

// keyPrefix 一般为业务前缀
var Cache Driver = NewMemoryCache("snake:", JsonEncoding{})

// 初始化缓存，在main.go里调用
// 默认是redis，这里也可以改为其他缓存，可以通过配置进行配置
func Init() {
	if gin.Mode() != gin.TestMode {
		// Cache = NewRedisCache()
		Cache = NewRedisCache(redis2.Client, "snake:", JsonEncoding{})
	}
}

type Driver interface {
	Set(key string, val interface{}, expiration time.Duration) error
	Get(key string) (interface{}, error)
	MultiSet(valMap map[string]interface{}, expiration time.Duration) error
	MultiGet(keys ...string) (interface{}, error)
	Del(keys ...string) error
	Incr(key string, step int64) (int64, error)
	Decr(key string, step int64) (int64, error)
}

func Set(key string, val interface{}, expiration time.Duration) error {
	return Cache.Set(key, val, expiration)
}

func Get(key string) (interface{}, error) {
	return Cache.Get(key)
}

func MultiSet(valMap map[string]interface{}, expiration time.Duration) error {
	return Cache.MultiSet(valMap, expiration)
}

func MultiGet(keys ...string) (interface{}, error) {
	return Cache.MultiGet(keys...)
}

func Del(keys ...string) error {
	return Cache.Del(keys...)
}

func Incr(key string, step int64) (int64, error) {
	return Cache.Incr(key, step)
}

func Decr(key string, step int64) (int64, error) {
	return Cache.Decr(key, step)
}
