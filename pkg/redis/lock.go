package redis

import (
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/go-redis/redis"
)

// Lock 定义lock结构体
type Lock struct {
	key         string
	redisClient *redis.Client
	timeout     time.Duration
}

// NewLock 实例化lock
func NewLock(conn *redis.Client, key string, defaultTimeout time.Duration) *Lock {
	return &Lock{
		key:         key,
		redisClient: conn,
		timeout:     defaultTimeout,
	}
}

// Lock 加锁
func (l *Lock) Lock(token string) (bool, error) {
	ok, err := l.redisClient.SetNX(l.GetKey(), token, l.timeout).Result()
	if err == redis.Nil {
		err = nil
	}
	return ok, err
}

// Unlock 解锁
// token 一致才会执行删除，避免误删，这里用了lua脚本进行事务处理
func (l *Lock) Unlock(token string) error {
	script := "if redis.call('get',KEYS[1]) == ARGV[1] then return redis.call('del',KEYS[1]) else return 0 end"
	_, err := l.redisClient.Eval(script, []string{l.GetKey()}, token).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetKey 获取key
func (l *Lock) GetKey() string {
	keyPrefix := viper.GetString("name")
	lockKey := "redis:lock"
	return strings.Join([]string{keyPrefix, lockKey, l.key}, ":")
}
