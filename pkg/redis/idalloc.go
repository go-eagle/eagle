package redis

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// IdAlloc id生成器
type IdAlloc struct {
	// key 为业务key, 比如生成用户id, 可以传入user_id
	key         string
	redisClient *redis.Client
}

// New 实例化
func New(conn *redis.Client, key string, defaultTimeout time.Duration) *Lock {
	return &Lock{
		key:         key,
		redisClient: conn,
		timeout:     defaultTimeout,
	}
}

// GetNewID 生成id
func (ia *IdAlloc) GetNewID(step int64) (int64, error) {
	key := ia.GetKey()
	id, err := ia.redisClient.IncrBy(key, step).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "redis incr err, key: %s", key)
	}
	return id, nil
}

// GetCurrentID 获取当前id
func (ia *IdAlloc) GetCurrentID() (int64, error) {
	key := ia.GetKey()
	ret, err := ia.redisClient.Get(key).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "redis get err, key: %s", key)
	}
	id, err := strconv.Atoi(ret)
	if err != nil {
		return 0, errors.Wrap(err, "str convert err")
	}
	return int64(id), nil
}

// GetKey 获取key, 由业务前缀+功能前缀+具体场景id组成
func (ia *IdAlloc) GetKey() string {
	keyPrefix := viper.GetString("name")
	lockKey := "idalloc"
	return strings.Join([]string{keyPrefix, lockKey, ia.key}, ":")
}
