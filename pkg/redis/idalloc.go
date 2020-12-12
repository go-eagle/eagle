package redis

import (
	"strconv"
	"strings"

	"github.com/1024casts/snake/pkg/conf"

	"github.com/1024casts/snake/pkg/log"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// IDAlloc id生成器
type IDAlloc struct {
	// key 为业务key, 由业务前缀+功能前缀+具体场景id组成
	// 比如生成用户id, 可以传入user_id， 完整示例: snake:idalloc:user_id
	key string
	// redis 实例，最好使用和业务独立的实例，最好可以部署集群，让 id alloc做到高可用
	redisClient *redis.Client
}

// New create a id alloc instance
func New(conn *redis.Client, key string) *IDAlloc {
	return &IDAlloc{
		key:         key,
		redisClient: conn,
	}
}

// GetNewID 生成id
func (ia *IDAlloc) GetNewID(step int64) (int64, error) {
	key := ia.GetKey()
	id, err := ia.redisClient.IncrBy(key, step).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "redis incr err, key: %s", key)
	}

	if id == 0 {
		log.Warnf("[redis.idalloc] %s GetNewID failed", ia.key)
		return 0, errors.Wrapf(err, "[redis.idalloc] %s GetNewID failed", ia.key)
	}
	return id, nil
}

// GetCurrentID 获取当前id
func (ia *IDAlloc) GetCurrentID() (int64, error) {
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

// GetKey 获取key
func (ia *IDAlloc) GetKey() string {
	keyPrefix := conf.Conf.App.Name
	lockKey := "idalloc"
	return strings.Join([]string{keyPrefix, lockKey, ia.key}, ":")
}
