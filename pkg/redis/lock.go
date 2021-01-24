package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

const (
	// LockKey redis lock key
	LockKey = "snake:redis:lock:%s"
	// DefaultTimeout default expire time
	DefaultTimeout = 2 * time.Second
)

type Option func(*Lock)

func Timeout(expiration time.Duration) Option {
	return func(l *Lock) {
		l.timeout = expiration
	}
}

// Lock 定义lock结构体
type Lock struct {
	key         string
	redisClient *redis.Client
	token       string
	timeout     time.Duration
	// todo: support retry
	// maxRetries int
}

// NewLock 实例化lock
func NewLock(conn *redis.Client, key string, options ...func(lock *Lock)) *Lock {
	lock := Lock{
		key:         key,
		redisClient: conn,
		timeout:     DefaultTimeout,
	}

	for _, option := range options {
		option(&lock)
	}
	return &lock
}

// Lock 加锁
func (l *Lock) Lock() (bool, error) {
	token := l.getToken()
	l.token = token
	ok, err := l.redisClient.SetNX(l.GetKey(), token, l.timeout).Result()
	if err == redis.Nil {
		err = nil
	}
	return ok, err
}

// Unlock 解锁
// token 一致才会执行删除，避免误删，这里用了lua脚本进行事务处理
func (l *Lock) Unlock() error {
	script := "if redis.call('get',KEYS[1]) == ARGV[1] then return redis.call('del',KEYS[1]) else return 0 end"
	_, err := l.redisClient.Eval(script, []string{l.GetKey()}, l.token).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetKey 获取key
func (l *Lock) GetKey() string {
	return fmt.Sprintf(LockKey, l.key)
}

// getToken 生成token
func (l *Lock) getToken() string {
	u, _ := uuid.NewRandom()
	return u.String()
}
