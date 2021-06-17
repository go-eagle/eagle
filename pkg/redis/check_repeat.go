package redis

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// PrefixCheckRepeat check repeat key
	PrefixCheckRepeat = "CHECK_REPEAT"
	// RepeatDefaultTimeout define default timeout
	RepeatDefaultTimeout = 60
)

// CheckRepeat define interface
type CheckRepeat interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	SetNX(key string, value interface{}, expiration time.Duration) (bool, error)
	Del(keys string) int64
}

type checkRepeat struct {
	client *redis.Client
}

// NewCheckRepeat create a check repeat
func NewCheckRepeat(client *redis.Client) CheckRepeat {
	return &checkRepeat{
		client: client,
	}
}

// GetKey 获取key
func getKey(key string) string {
	return strings.Join([]string{PrefixCheckRepeat, key}, ":")
}

func (c *checkRepeat) Set(key string, value interface{}, expiration time.Duration) error {
	key = getKey(key)
	return c.client.Set(context.Background(), key, value, expiration).Err()
}

func (c *checkRepeat) Get(key string) (string, error) {
	key = getKey(key)
	return c.client.Get(context.Background(), key).Result()
}

func (c *checkRepeat) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	key = getKey(key)
	return c.client.SetNX(context.Background(), key, value, expiration).Result()
}

func (c *checkRepeat) Del(key string) int64 {
	key = getKey(key)
	var keys []string
	keys = append(keys, key)
	return c.client.Del(context.Background(), keys...).Val()
}
