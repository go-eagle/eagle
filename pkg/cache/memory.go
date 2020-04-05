package cache

import (
	"bytes"
	"encoding/gob"
	"sync"
	"time"

	"github.com/lexkong/log"

	"github.com/pkg/errors"
)

type memoryCache struct {
	client    *sync.Map
	KeyPrefix string
	encoding  Encoding
}

func NewMemoryCache(keyPrefix string, encoding Encoding) *memoryCache {
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

// interface 转 byte
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

func (m memoryCache) Set(key string, val interface{}, expiration time.Duration) error {
	cacheKey, err := BuildCacheKey(m.KeyPrefix, key)
	if err != nil {
		return errors.Wrapf(err, "build cache key err, key is %+v", key)
	}
	m.client.Store(cacheKey, newItem(val, expiration))
	return nil
}

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

func (m memoryCache) MultiSet(valMap map[string]interface{}, expiration time.Duration) error {
	panic("implement me")
}

func (m memoryCache) MultiGet(keys ...string) (interface{}, error) {
	panic("implement me")
}

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

func (m memoryCache) Incr(key string, step int64) (int64, error) {
	panic("implement me")
}

func (m memoryCache) Decr(key string, step int64) (int64, error) {
	panic("implement me")
}
