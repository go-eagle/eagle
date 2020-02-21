package cache

import (
	"time"

	redis2 "github.com/1024casts/snake/pkg/redis"
	"github.com/go-redis/redis"
)

type redisCache struct {
	client *redis.Client
}

func newRedisCache() *redisCache {
	return &redisCache{
		client: redis2.Client,
	}
}

func (c *redisCache) Set(key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(key, value, expiration).Err()
}

func (c *redisCache) Get(key string) ([]byte, error) {
	return c.client.Get(key).Bytes()
}

func (c *redisCache) Del(keys ...string) (int64, error) {
	return c.client.Del(keys...).Result()
}

func (c *redisCache) Incr(key string, step int64) (int64, error) {
	return c.client.IncrBy(key, step).Result()
}

func (c *redisCache) Decr(key string, step int64) (int64, error) {
	return c.client.DecrBy(key, step).Result()
}
