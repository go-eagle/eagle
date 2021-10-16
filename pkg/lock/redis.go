package lock

import (
	"context"
	"fmt"
	"time"

	"github.com/go-eagle/eagle/pkg/log"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// redisLock is a redis lock.
type redisLock struct {
	key         string
	redisClient *redis.Client
	token       string
	timeout     time.Duration
}

// Option .
type Option func(*redisLock)

// WithTimeout with a timeout
func WithTimeout(expiration time.Duration) Option {
	return func(l *redisLock) {
		l.timeout = expiration
	}
}

// NewRedisLock new a redis lock instance
// nolint
func NewRedisLock(rdb *redis.Client, key string, opts ...Option) *redisLock {
	lock := redisLock{
		key:         key,
		redisClient: rdb,
		token:       genToken(),
		timeout:     DefaultTimeout,
	}

	for _, o := range opts {
		o(&lock)
	}
	return &lock
}

// Lock acquires the lock.
func (l *redisLock) Lock(ctx context.Context) (bool, error) {
	isSet, err := l.redisClient.SetNX(ctx, l.GetKey(), l.token, l.timeout).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		log.Errorf("acquires the lock err, key: %s, err: %s", l.GetKey(), err.Error())
		return false, err
	}
	return isSet, nil
}

// Unlock del the lock.
// NOTE: token 一致才会执行删除，避免误删，这里用了lua脚本进行事务处理
func (l *redisLock) Unlock(ctx context.Context) (bool, error) {
	luaScript := "if redis.call('GET',KEYS[1]) == ARGV[1] then return redis.call('DEL',KEYS[1]) else return 0 end"
	ret, err := l.redisClient.Eval(ctx, luaScript, []string{l.GetKey()}, l.token).Result()
	if err != nil {
		return false, err
	}
	reply, ok := ret.(int64)
	if !ok {
		return false, nil
	}
	return reply == 1, nil
}

// GetKey 获取key
func (l *redisLock) GetKey() string {
	return fmt.Sprintf(LockKey, l.key)
}

// genToken 生成token
func genToken() string {
	u, _ := uuid.NewRandom()
	return u.String()
}
