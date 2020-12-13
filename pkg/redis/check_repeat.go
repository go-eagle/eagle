package redis

import (
	"strings"
	"time"

	"github.com/1024casts/snake/pkg/conf"

	"github.com/go-redis/redis"
)

const (
	PrefixCheckRepeat    = "CHECK_REPEAT"
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
	keyPrefix := conf.Conf.App.Name
	return strings.Join([]string{keyPrefix, PrefixCheckRepeat, key}, ":")
}

func (c *checkRepeat) Set(key string, value interface{}, expiration time.Duration) error {
	key = getKey(key)
	return c.client.Set(key, value, expiration).Err()
}

func (c *checkRepeat) Get(key string) (string, error) {
	key = getKey(key)
	return c.client.Get(key).Result()
}

func (c *checkRepeat) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	key = getKey(key)
	return c.client.SetNX(key, value, expiration).Result()
}

func (c *checkRepeat) Del(key string) int64 {
	key = getKey(key)
	var keys []string
	keys = append(keys, key)
	return c.client.Del(keys...).Val()
}
