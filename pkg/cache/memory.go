package cache

import (
	"reflect"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/1024casts/snake/pkg/log"
)

type memoryCache struct {
	Store     *sync.Map
	KeyPrefix string
	encoding  Encoding
}

// NewMemoryCache 实例化一个内存cache
func NewMemoryCache(keyPrefix string, encoding Encoding) Driver {
	return &memoryCache{
		Store:     &sync.Map{},
		KeyPrefix: keyPrefix,
		encoding:  encoding,
	}
}

// item 存储的对象
type itemWithTTL struct {
	expires int64
	value   []byte
}

// newItem 返回带有效期的value
func newItem(value []byte, expires time.Duration) itemWithTTL {
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
func getValue(item interface{}, ok bool) ([]byte, bool) {
	if !ok {
		return nil, false
	}

	var itemObj itemWithTTL
	if itemObj, ok = item.(itemWithTTL); !ok {
		return nil, false
	}

	// 过期返回空
	if itemObj.expires > 0 && itemObj.expires < time.Now().Unix() {
		return nil, false
	}

	return itemObj.value, true
}

// GarbageCollect 回收已过期的缓存，可以通过定时任务来触发
func (m *memoryCache) GarbageCollect() {
	m.Store.Range(func(key, value interface{}) bool {
		if item, ok := value.(itemWithTTL); ok {
			if item.expires > 0 && item.expires < time.Now().Unix() {
				log.Info("回收垃圾[%s]", key.(string))
				m.Store.Delete(key)
			}
		}
		return true
	})
}

// Set data
func (m *memoryCache) Set(key string, val interface{}, expiration time.Duration) error {
	buf, err := Marshal(m.encoding, val)
	if err != nil {
		return errors.Wrapf(err, "marshal data err, value is %+v", val)
	}
	cacheKey, err := BuildCacheKey(m.KeyPrefix, key)
	if err != nil {
		return errors.Wrapf(err, "build cache key err, key is %+v", key)
	}
	m.Store.Store(cacheKey, newItem(buf, expiration))
	return nil
}

// Get data
func (m *memoryCache) Get(key string, val interface{}) error {
	cacheKey, err := BuildCacheKey(m.KeyPrefix, key)
	if err != nil {
		return errors.Wrapf(err, "build cache key err, key is %+v", key)
	}
	data, ok := getValue(m.Store.Load(cacheKey))
	if !ok {
		return nil
	}
	err = Unmarshal(m.encoding, data, val)
	if err != nil {
		return errors.Wrapf(err, "unmarshal data error, key=%s, cacheKey=%s type=%v, json is %+v ",
			key, cacheKey, reflect.TypeOf(val), string(data))
	}
	return nil
}

// MultiSet 批量set
func (m *memoryCache) MultiSet(valMap map[string]interface{}, expiration time.Duration) error {
	panic("implement me")
}

// MultiGet 批量获取
func (m *memoryCache) MultiGet(keys []string, val interface{}) error {
	panic("implement me")
}

// Del 批量删除
func (m *memoryCache) Del(keys ...string) error {
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
		m.Store.Delete(cacheKey)
	}
	return nil
}

// Incr 自增
func (m *memoryCache) Incr(key string, step int64) (int64, error) {
	panic("implement me")
}

// Decr 自减
func (m *memoryCache) Decr(key string, step int64) (int64, error) {
	panic("implement me")
}
