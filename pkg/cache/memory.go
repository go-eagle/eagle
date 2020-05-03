package cache

import (
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/1024casts/snake/pkg/log"
)

type memoryCache struct {
	client    *sync.Map
	KeyPrefix string
	encoding  Encoding
}

// NewMemoryCache 实例化一个内存cache
func NewMemoryCache(keyPrefix string, encoding Encoding) Driver {
	return &memoryCache{
		client:    &sync.Map{},
		KeyPrefix: keyPrefix,
		encoding:  encoding,
	}
}

// item 存储的对象
type itemWithTTL struct {
	expires int64
	value   interface{}
}

// newItem 返回带有效期的value
func newItem(value interface{}, expires time.Duration) itemWithTTL {
	expires64 := int64(expires)
	if expires > 0 {
		expires64 = time.Now().Unix() + expires64
	}
	return itemWithTTL{
		value:   value,
		expires: expires64,
	}
}

// getValue 从itemWithTTL中取值
func getValue(item interface{}, ok bool) (interface{}, bool) {
	if !ok {
		return nil, false
	}

	var itemObj itemWithTTL
	if itemObj, ok = item.(itemWithTTL); !ok {
		return nil, false
	}

	if itemObj.expires > 0 && itemObj.expires < time.Now().Unix() {
		return nil, false
	}

	return itemObj.value, true
}

// Set data
func (m memoryCache) Set(key string, val interface{}, expiration time.Duration) error {
	cacheKey, err := BuildCacheKey(m.KeyPrefix, key)
	if err != nil {
		return errors.Wrapf(err, "build cache key err, key is %+v", key)
	}
	m.client.Store(cacheKey, newItem(val, expiration))
	return nil
}

// Get data
func (m memoryCache) Get(key string) (interface{}, error) {
	cacheKey, err := BuildCacheKey(m.KeyPrefix, key)
	if err != nil {
		return nil, errors.Wrapf(err, "build cache key err, key is %+v", key)
	}
	val, ok := getValue(m.client.Load(cacheKey))
	if !ok {
		return nil, errors.New("memory get value err")
	}
	return val, nil
}

// MultiSet 批量set
func (m memoryCache) MultiSet(valMap map[string]interface{}, expiration time.Duration) error {
	panic("implement me")
}

// MultiGet 批量获取
func (m memoryCache) MultiGet(keys ...string) (interface{}, error) {
	panic("implement me")
}

// Del 批量删除
func (m memoryCache) Del(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	// 批量构建cacheKey
	for _, key := range keys {
		cacheKey, err := BuildCacheKey(m.KeyPrefix, key)
		if err != nil {
			log.Warnf("build cache key err: %+v, key is %+v", err, key)
			continue
		}
		m.client.Delete(cacheKey)
	}
	return nil
}

// Incr 自增
func (m memoryCache) Incr(key string, step int64) (int64, error) {
	panic("implement me")
}

// Decr 自减
func (m memoryCache) Decr(key string, step int64) (int64, error) {
	panic("implement me")
}
