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
	return nil
}

func (c *redisCache) Get(key string) (interface{}, error) {
	return nil, nil
}

func (c *redisCache) Del(string) error {
	panic("implement me")
}
