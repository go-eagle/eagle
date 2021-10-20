package lock

import (
	"context"
	"fmt"
	"time"

	"github.com/go-eagle/eagle/pkg/log"

	"github.com/go-redis/redis/v8"
)

// RedisLock is a redis lock.
type RedisLock struct {
	key         string
	redisClient *redis.Client
	token       string
}

// NewRedisLock new a redis lock instance
// nolint
func NewRedisLock(rdb *redis.Client, key string) *RedisLock {
	opt := &RedisLock{
		key:         getRedisKey(key),
		redisClient: rdb,
		token:       genToken(),
	}
	return opt
}

// Lock acquires the lock.
func (l *RedisLock) Lock(ctx context.Context, timeout time.Duration) (bool, error) {
	isSet, err := l.redisClient.SetNX(ctx, l.key, l.token, timeout).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		log.Errorf("acquires the lock err, key: %s, err: %s", l.key, err.Error())
		return false, err
	}
	return isSet, nil
}

// Unlock del the lock.
// NOTE: token 一致才会执行删除，避免误删，这里用了lua脚本进行事务处理
func (l *RedisLock) Unlock(ctx context.Context) (bool, error) {
	luaScript := "if redis.call('GET',KEYS[1]) == ARGV[1] then return redis.call('DEL',KEYS[1]) else return 0 end"
	ret, err := l.redisClient.Eval(ctx, luaScript, []string{l.key}, l.token).Result()
	if err != nil {
		return false, err
	}
	reply, ok := ret.(int64)
	if !ok {
		return false, nil
	}
	return reply == 1, nil
}

// getRedisKey 获取key
func getRedisKey(key string) string {
	return fmt.Sprintf(RedisLockKey, key)
}
