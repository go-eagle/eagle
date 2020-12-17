// Package pkg 计数器，可以用于业务的各种模型统计使用
// 场景：常用语重复策略，或者反作弊处理控制
package pkg

import (
	"strings"
	"time"

	"github.com/go-redis/redis"

	"github.com/1024casts/snake/pkg/conf"
	redis2 "github.com/1024casts/snake/pkg/redis"
)

const (
	// PrefixCounter counter key
	PrefixCounter = "COUNTER"
	// DefaultStep default step key
	DefaultStep           = 1
	DefaultExpirationTime = 600 * time.Second
)

// Counter define struct
type Counter struct {
	client *redis.Client
}

// NewCounter create a counter
func NewCounter() *Counter {
	return &Counter{
		client: redis2.RedisClient,
	}
}

// GetKey 获取key
func (c *Counter) GetKey(key string) string {
	keyPrefix := conf.Conf.App.Name
	return strings.Join([]string{keyPrefix, PrefixCounter, key}, ":")
}

// SetCounter set counter
func (c *Counter) SetCounter(idStr string, expiration time.Duration) (int64, error) {
	key := c.GetKey(idStr)
	ret, err := c.client.IncrBy(key, DefaultStep).Result()
	if err != nil {
		return 0, err
	}
	_, _ = c.client.Expire(key, expiration).Result()
	return ret, nil
}

// GetCounter get total count
func (c *Counter) GetCounter(idStr string) (int64, error) {
	key := c.GetKey(idStr)
	return c.client.Get(key).Int64()
}

// DelCounter del count
func (c *Counter) DelCounter(idStr string) int64 {
	key := c.GetKey(idStr)
	var keys []string
	keys = append(keys, key)
	return c.client.Del(keys...).Val()
}
