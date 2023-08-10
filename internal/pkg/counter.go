// Package pkg 计数器，可以用于业务的各种模型统计使用
// 场景：常用于重复策略，或者反作弊处理控制
package pkg

import (
	"context"
	"fmt"
	"time"

	redis2 "github.com/go-eagle/eagle/pkg/redis"
	"github.com/redis/go-redis/v9"
)

const (
	// PrefixCounter counter key
	PrefixCounter = "eagle:counter:%s"
	// DefaultStep default step key
	DefaultStep = 1
	// DefaultExpirationTime .
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
	return fmt.Sprintf(PrefixCounter, key)
}

// SetCounter set counter
func (c *Counter) SetCounter(ctx context.Context, idStr string, expiration time.Duration) (int64, error) {
	key := c.GetKey(idStr)
	ret, err := c.client.IncrBy(ctx, key, DefaultStep).Result()
	if err != nil {
		return 0, err
	}
	_, _ = c.client.Expire(ctx, key, expiration).Result()
	return ret, nil
}

// GetCounter get total count
func (c *Counter) GetCounter(ctx context.Context, idStr string) (int64, error) {
	key := c.GetKey(idStr)
	return c.client.Get(ctx, key).Int64()
}

// DelCounter del count
func (c *Counter) DelCounter(ctx context.Context, idStr string) int64 {
	key := c.GetKey(idStr)
	var keys []string
	keys = append(keys, key)
	return c.client.Del(ctx, keys...).Val()
}
