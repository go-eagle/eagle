package cache

import (
	"time"
)

type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) ([]byte, error)
	Del(keys ...string) (int64, error)
	Incr(key string, step int64) (int64, error)
	Decr(key string, step int64) (int64, error)
}
