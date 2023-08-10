package lock

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-eagle/eagle/pkg/redis"
)

func TestLockWithDefaultTimeout(t *testing.T) {
	redis.InitTestRedis()
	expiration := 2 * time.Second

	lock := NewRedisLock(redis.RedisClient, "lock1", expiration)
	ok, err := lock.Lock(context.Background())
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Fatal("lock is not ok")
	}

	ok, err = lock.Unlock(context.Background())
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("Unlock is not ok")
	}

	t.Log(ok)
}

func TestLockWithTimeout(t *testing.T) {
	redis.InitTestRedis()
	expiration := 2 * time.Second

	t.Run("should lock/unlock success", func(t *testing.T) {
		ctx := context.Background()
		lock1 := NewRedisLock(redis.RedisClient, "lock2", expiration)
		ok, err := lock1.Lock(ctx)
		assert.Nil(t, err)
		assert.True(t, ok)

		ok, err = lock1.Unlock(ctx)
		assert.Nil(t, err)
		assert.True(t, ok)
	})

	t.Run("should unlock failed", func(t *testing.T) {
		ctx := context.Background()
		lock2 := NewRedisLock(redis.RedisClient, "lock3", expiration)
		ok, err := lock2.Lock(ctx)
		assert.Nil(t, err)
		assert.True(t, ok)

		time.Sleep(3 * time.Second)

		ok, err = lock2.Unlock(ctx)
		fmt.Println("===*****************", ok, err)
		assert.Nil(t, err)
		assert.True(t, ok)
	})
}
