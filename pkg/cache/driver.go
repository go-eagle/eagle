package cache

import (
	"time"

	"github.com/spf13/viper"

	redis2 "github.com/1024casts/snake/pkg/redis"
)

// Cache 生成一个缓存客户端，其中keyPrefix 一般为业务前缀
var Cache Driver

const (
	// memCacheDriver 内存缓存
	memCacheDriver = "memory"
	// redisCacheDriver redis缓存
	redisCacheDriver = "redis"
)

// Init 初始化缓存，在main.go里调用
func Init() {
	cacheDriver := viper.GetString("cache.driver")
	cachePrefix := viper.GetString("cache.prefix")
	encoding := JSONEncoding{}

	switch cacheDriver {
	case memCacheDriver:
		Cache = NewMemoryCache(cachePrefix, encoding)
	case redisCacheDriver:
		Cache = NewRedisCache(redis2.Client, cachePrefix, encoding)
	default:
		Cache = NewMemoryCache(cachePrefix, encoding)
	}
}

// Driver 定义cache驱动接口
type Driver interface {
	Set(key string, val interface{}, expiration time.Duration) error
	Get(key string) (interface{}, error)
	MultiSet(valMap map[string]interface{}, expiration time.Duration) error
	MultiGet(keys ...string) (interface{}, error)
	Del(keys ...string) error
	Incr(key string, step int64) (int64, error)
	Decr(key string, step int64) (int64, error)
}

// Set 数据
func Set(key string, val interface{}, expiration time.Duration) error {
	return Cache.Set(key, val, expiration)
}

// Get 数据
func Get(key string) (interface{}, error) {
	return Cache.Get(key)
}

// MultiSet 批量set
func MultiSet(valMap map[string]interface{}, expiration time.Duration) error {
	return Cache.MultiSet(valMap, expiration)
}

// MultiGet 批量获取
func MultiGet(keys ...string) (interface{}, error) {
	return Cache.MultiGet(keys...)
}

// Del 批量删除
func Del(keys ...string) error {
	return Cache.Del(keys...)
}

// Incr 自增
func Incr(key string, step int64) (int64, error) {
	return Cache.Incr(key, step)
}

// Decr 自减
func Decr(key string, step int64) (int64, error) {
	return Cache.Decr(key, step)
}
