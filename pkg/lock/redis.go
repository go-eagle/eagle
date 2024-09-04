package lock

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/go-eagle/eagle/pkg/log"
)

const (
	// renewalDuration is the renewal duration.
	renewalDuration int64 = 1000
)

// Lua scripts for locking and unlocking
// It only init once and cache in memory
var (

	// lockScript init lua script
	lockScript = redis.NewScript(lockLuaScript)

	// unlockScript init lua script
	unlockScript = redis.NewScript(unlockLuaScript)
)

// RedisLock is a Redis-based distributed lock.
type RedisLock struct {
	key         string
	redisClient *redis.Client
	token       string
	expiration  time.Duration
	mu          sync.Mutex    // 用于保护共享属性
	renewing    bool          // 续期标志
	stopRenew   chan struct{} // 用于停止续期
}

// NewRedisLock new a redis lock instance
// nolint
func NewRedisLock(rdb *redis.Client, key string, expiration time.Duration) *RedisLock {
	opt := &RedisLock{
		key:         getRedisKey(key),
		redisClient: rdb,
		token:       genToken(),
		expiration:  expiration,
		stopRenew:   make(chan struct{}),
	}
	return opt
}

// Lock acquires the lock.
// It will return false if the lock is already acquired.
func (l *RedisLock) Lock(ctx context.Context) (bool, error) {
	// 加锁，防止并发问题
	l.mu.Lock()
	defer l.mu.Unlock()

	ret, err := lockScript.Run(ctx, l.redisClient, []string{l.key},
		[]string{l.token, strconv.FormatInt(l.expiration.Milliseconds()+renewalDuration, 10)},
	).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		log.WithContext(ctx).Errorf("redis lock: acquires the lock err, key: %s, err: %s", l.key, err.Error())
		return false, err
	} else if ret == nil {
		return false, nil
	}

	reply, ok := ret.(string)
	if ok && reply == "OK" {
		if !l.renewing {
			l.renewing = true
			go l.autoRenew(ctx) // 启动续期协程
		}
		return true, nil
	}

	return false, nil
}

// Unlock release a lock.
// NOTE: token 一致才会执行删除，避免误删，这里用了lua脚本进行事务处理
func (l *RedisLock) Unlock(ctx context.Context) (bool, error) {
	// 解锁时也需要加锁保护
	l.mu.Lock()
	defer l.mu.Unlock()

	// 停止续期协程并标记不再续期
	if l.renewing {
		close(l.stopRenew)
		l.renewing = false
	}

	for i := 0; i < 3; i++ { // 最多重试3次
		ret, err := unlockScript.Run(ctx, l.redisClient, []string{l.key}, l.token).Result()
		if err != nil {
			log.WithContext(ctx).Errorf("redis lock: failed to unlock, attempt %d, key: %s, err: %v", i+1, l.key, err)
			time.Sleep(50 * time.Millisecond) // 等待一下再重试
			continue
		}
		reply, ok := ret.(int64)
		if ok && reply == 1 {
			return true, nil
		}
		break
	}

	return false, errors.New("redis lock: failed to unlock after multiple attempts")
}

func (l *RedisLock) autoRenew(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(renewalDuration) * time.Millisecond / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 续期操作
			l.mu.Lock()
			l.redisClient.Expire(ctx, l.key, l.expiration)
			l.mu.Unlock()
		case <-l.stopRenew:
			return
		case <-ctx.Done():
			return
		}
	}
}

// GetTTL returns the TTL of the lock.
func (l *RedisLock) GetTTL(ctx context.Context) (time.Duration, error) {
	return l.redisClient.TTL(ctx, l.key).Result()
}

// getRedisKey returns the Redis key for the lock.
func getRedisKey(key string) string {
	return fmt.Sprintf(RedisLockKey, key)
}
