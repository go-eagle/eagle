package lock

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-eagle/eagle/pkg/redis"
)

func TestLockWithDefaultTimeout(t *testing.T) {
	redis.InitTestRedis()

	lock := NewRedisLock(redis.RedisClient, "lock1")
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

	t.Run("should lock/unlock is success", func(t *testing.T) {
		lock1 := NewRedisLock(redis.RedisClient, "lock2", WithTimeout(2*time.Second))
		ok, err := lock1.Lock(context.Background())
		assert.Nil(t, err)
		assert.True(t, ok)

		ok, err = lock1.Unlock(context.Background())
		assert.Nil(t, err)
		assert.True(t, ok)
	})

	t.Run("should unlock is failed", func(t *testing.T) {
		lock2 := NewRedisLock(redis.RedisClient, "lock3", WithTimeout(2*time.Second))
		ok, err := lock2.Lock(context.Background())
		assert.Nil(t, err)
		assert.True(t, ok)

		time.Sleep(3 * time.Second)

		ok, err = lock2.Unlock(context.Background())
		assert.Nil(t, err)
		assert.True(t, ok)
	})

}
