package cache

import (
	"reflect"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"

	"github.com/1024casts/snake/pkg/encoding"
	"github.com/1024casts/snake/pkg/log"
)

// redisCache redis cache结构体
type redisCache struct {
	client            *redis.Client
	KeyPrefix         string
	encoding          encoding.Encoding
	DefaultExpireTime time.Duration
	newObject         func() interface{}
}

// NewRedisCache new一个cache, client 参数是可传入的，方便进行单元测试
func NewRedisCache(client *redis.Client, keyPrefix string, encoding encoding.Encoding, newObject func() interface{}) Cache {
	return &redisCache{
		client:    client,
		KeyPrefix: keyPrefix,
		encoding:  encoding,
		newObject: newObject,
	}
}

func (c *redisCache) Set(key string, val interface{}, expiration time.Duration) error {
	buf, err := encoding.Marshal(c.encoding, val)
	if err != nil {
		return errors.Wrapf(err, "marshal data err, value is %+v", val)
	}

	cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
	if err != nil {
		return errors.Wrapf(err, "build cache key err, key is %+v", key)
	}
	if expiration == 0 {
		expiration = DefaultExpireTime
	}
	err = c.client.Set(cacheKey, buf, expiration).Err()
	if err != nil {
		return errors.Wrapf(err, "redis set error")
	}
	return nil
}

func (c *redisCache) Get(key string, val interface{}) error {
	cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
	if err != nil {
		return errors.Wrapf(err, "build cache key err, key is %+v", key)
	}

	data, err := c.client.Get(cacheKey).Bytes()
	if err != nil {
		if err != redis.Nil {
			return errors.Wrapf(err, "get data error from redis, key is %+v", cacheKey)
		}
	}

	// 防止data为空时，Unmarshal报错
	if string(data) == "" {
		return nil
	}
	if string(data) == NotFoundPlaceholder {
		return ErrPlaceholder
	}
	err = encoding.Unmarshal(c.encoding, data, val)
	if err != nil {
		return errors.Wrapf(err, "unmarshal data error, key=%s, cacheKey=%s type=%v, json is %+v ",
			key, cacheKey, reflect.TypeOf(val), string(data))
	}
	return nil
}

func (c *redisCache) MultiSet(valueMap map[string]interface{}, expiration time.Duration) error {
	if len(valueMap) == 0 {
		return nil
	}
	// key-value是成对的，所以这里的容量是map的2倍
	paris := make([]interface{}, 0, 2*len(valueMap))
	for key, value := range valueMap {
		buf, err := encoding.Marshal(c.encoding, value)
		if err != nil {
			log.Warnf("marshal data err: %+v, value is %+v", err, value)
			continue
		}
		cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
		if err != nil {
			log.Warnf("build cache key err: %+v, key is %+v", err, key)
			continue
		}
		if expiration == 0 {
			expiration = DefaultExpireTime
		}
		paris = append(paris, []byte(cacheKey))
		paris = append(paris, buf)
	}
	if expiration == 0 {
		expiration = DefaultExpireTime
	}
	err := c.client.MSet(paris...).Err()
	if err != nil {
		return errors.Wrapf(err, "redis multi set error")
	}
	for i := 0; i < len(paris); i = i + 2 {
		switch paris[i].(type) {
		case []byte:
			c.client.Expire(string(paris[i].([]byte)), expiration)
		default:
			log.Warnf("redis expire is unsupported key type: %+v", reflect.TypeOf(paris[i]))
		}
	}
	return nil
}

func (c *redisCache) MultiGet(keys []string, value interface{}) error {
	if len(keys) == 0 {
		return nil
	}
	cacheKeys := make([]string, len(keys))
	for index, key := range keys {
		cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
		if err != nil {
			return errors.Wrapf(err, "build cache key err, key is %+v", key)
		}
		cacheKeys[index] = cacheKey
	}
	values, err := c.client.MGet(cacheKeys...).Result()
	if err != nil {
		return errors.Wrapf(err, "redis MGet error, keys is %+v", keys)
	}

	// 通过反射注入到map
	valueMap := reflect.ValueOf(value)
	for i, value := range values {
		if value == nil {
			continue
		}
		object := c.newObject()
		err = encoding.Unmarshal(c.encoding, []byte(value.(string)), object)
		if err != nil {
			log.Warnf("unmarshal data error: %+v, key=%s, cacheKey=%s type=%v", err,
				keys[i], cacheKeys[i], reflect.TypeOf(value))
			continue
		}
		valueMap.SetMapIndex(reflect.ValueOf(cacheKeys[i]), reflect.ValueOf(object))
	}
	return nil
}

func (c *redisCache) Del(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	// 批量构建cacheKey
	cacheKeys := make([]string, len(keys))
	for index, key := range keys {
		cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
		if err != nil {
			log.Warnf("build cache key err: %+v, key is %+v", err, key)
			continue
		}
		cacheKeys[index] = cacheKey
	}
	err := c.client.Del(cacheKeys...).Err()
	if err != nil {
		return errors.Wrapf(err, "redis delete error, keys is %+v", keys)
	}
	return nil
}

func (c *redisCache) SetCacheWithNotFound(key string) error {
	return c.client.Set(key, NotFoundPlaceholder, DefaultNotFoundExpireTime).Err()
}
