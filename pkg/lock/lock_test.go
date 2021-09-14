package lock

import (
	"context"
	"github.com/go-eagle/eagle/pkg/redis"
	"testing"
	"time"
)

func TestLockWithDefaultTimeout(t *testing.T) {
	redis.InitTestRedis()

	lock := NewLock(redis.RedisClient, "test:lock")
	ok, err := lock.Lock(context.Background())
	if err != nil {
		t.Error(err)
	}

	err = lock.Unlock(context.Background())
	if err != nil {
		t.Error(err)
	}

	t.Log(ok)
}

func TestLockWithTimeout(t *testing.T) {
	redis.InitTestRedis()
	lock := NewLock(redis.RedisClient, "test:lock", WithTimeout(3*time.Second))
	ok, err := lock.Lock(context.Background())
	if err != nil {
		t.Error(err)
	}

	err = lock.Unlock(context.Background())
	if err != nil {
		t.Error(err)
	}

	t.Log(ok)
}
