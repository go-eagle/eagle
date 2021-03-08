package redis

import (
	"strconv"
	"strings"

	"github.com/1024casts/snake/config"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"

	"github.com/1024casts/snake/pkg/log"
)

// IDAlloc id生成器
// key 为业务key, 由业务前缀+功能前缀+具体场景id组成
// 比如生成用户id, 可以传入user_id， 完整示例: snake:idalloc:user_id
type IDAlloc struct {
	// redis 实例，最好使用和业务独立的实例，最好可以部署集群，让 id alloc做到高可用
	redisClient *redis.Client
}

// NewIDAlloc create a id alloc instance
func NewIDAlloc(conn *redis.Client) *IDAlloc {
	return &IDAlloc{
		redisClient: conn,
	}
}

// GetNewID 生成id
func (ia *IDAlloc) GetNewID(key string, step int64) (int64, error) {
	key = ia.GetKey(key)
	id, err := ia.redisClient.IncrBy(key, step).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "redis incr err, key: %s", key)
	}

	if id == 0 {
		log.Warnf("[redis.idalloc] %s GetNewID failed", key)
		return 0, errors.Wrapf(err, "[redis.idalloc] %s GetNewID failed", key)
	}
	return id, nil
}

// GetCurrentID 获取当前id
func (ia *IDAlloc) GetCurrentID(key string) (int64, error) {
	key = ia.GetKey(key)
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
func (ia *IDAlloc) GetKey(key string) string {
	keyPrefix := config.Conf.App.Name
	lockKey := "idalloc"
	return strings.Join([]string{keyPrefix, lockKey, key}, ":")
}
