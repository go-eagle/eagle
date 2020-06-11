package cache

import (
	"fmt"
	"time"

	"github.com/1024casts/snake/pkg/redis"
)

// Client 生成一个缓存客户端，其中keyPrefix 一般为业务前缀
var Client Driver

const (
	// memCacheDriver 内存缓存
	memCacheDriver = "memory"
	// redisCacheDriver redis缓存
	redisCacheDriver = "redis"
)

// Init 初始化缓存，在main.go里调用
func Init() {
	// todo: 读取配置文件
	cacheDriver := "redis"
	cachePrefix := "snake"
	fmt.Println("get prefix key1:", cachePrefix)
	encoding := JSONEncoding{}

	switch cacheDriver {
	case memCacheDriver:
		Client = NewMemoryCache(cachePrefix, encoding)
	case redisCacheDriver:
		// todo: redis.Init() 已经在 main.go 执行，这里应该不需要再初始化，待排查
		redis.Init()
		Client = NewRedisCache(redis.Client, cachePrefix, encoding)
	default:
		Client = NewMemoryCache(cachePrefix, encoding)
	}
}

// Driver 定义cache驱动接口
type Driver interface {
	Set(key string, val interface{}, expiration time.Duration) error
	Get(key string, val interface{}) error
	MultiSet(valMap map[string]interface{}, expiration time.Duration) error
	MultiGet(keys []string, val interface{}) error
	Del(keys ...string) error
	Incr(key string, step int64) (int64, error)
	Decr(key string, step int64) (int64, error)
}

// Set 数据
func Set(key string, val interface{}, expiration time.Duration) error {
	return Client.Set(key, val, expiration)
}

// Get 数据
func Get(key string, val interface{}) error {
	return Client.Get(key, val)
}

// MultiSet 批量set
func MultiSet(valMap map[string]interface{}, expiration time.Duration) error {
	return Client.MultiSet(valMap, expiration)
}

// MultiGet 批量获取
func MultiGet(keys []string, val interface{}) error {
	return Client.MultiGet(keys, val)
}

// Del 批量删除
func Del(keys ...string) error {
	return Client.Del(keys...)
}

// Incr 自增
func Incr(key string, step int64) (int64, error) {
	return Client.Incr(key, step)
}

// Decr 自减
func Decr(key string, step int64) (int64, error) {
	return Client.Decr(key, step)
}
