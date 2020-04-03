package cache

import (
	"time"

	"github.com/gin-gonic/gin"
)

var Cache Driver = NewMemoryCache()

// 初始化缓存
func Init() {
	if gin.Mode() != gin.TestMode {
		Cache = NewRedisCache()
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
