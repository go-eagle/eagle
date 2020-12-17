package redis

import (
	"strings"
	"time"

	"github.com/1024casts/snake/pkg/conf"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

const (
	// LockKey redis lock key
	LockKey = "redis:lock"
	// DefaultTimeout default expire time
	DefaultTimeout = 2 * time.Second
)

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
func NewLock(conn *redis.Client, key string) *Lock {
	return &Lock{
		key:         key,
		redisClient: conn,
		timeout:     DefaultTimeout,
	}
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

// SetExpireTime set timeout time
func (l *Lock) SetExpireTime(expiration time.Duration) {
	_, _ = l.redisClient.Expire(l.key, expiration).Result()
}

// GetKey 获取key
func (l *Lock) GetKey() string {
	keyPrefix := conf.Conf.App.Name
	return strings.Join([]string{keyPrefix, LockKey, l.key}, ":")
}

// getToken 生成token
func (l *Lock) getToken() string {
	u, _ := uuid.NewRandom()
	return u.String()
}
