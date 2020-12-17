package cache

import (
	"reflect"
	"time"

	"github.com/dgraph-io/ristretto"

	"github.com/pkg/errors"

	"github.com/1024casts/snake/pkg/log"
)

type memoryCache struct {
	Store     *ristretto.Cache
	KeyPrefix string
	encoding  Encoding
}

// NewMemoryCache 实例化一个内存cache
func NewMemoryCache(keyPrefix string, encoding Encoding) Driver {
	config := &ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	}
	store, _ := ristretto.NewCache(config)
	return &memoryCache{
		Store:     store,
		KeyPrefix: keyPrefix,
		encoding:  encoding,
	}
}

// Set add cache
func (m *memoryCache) Set(key string, val interface{}, expiration time.Duration) error {
	buf, err := Marshal(m.encoding, val)
	if err != nil {
		return errors.Wrapf(err, "marshal data err, value is %+v", val)
	}
	cacheKey, err := BuildCacheKey(m.KeyPrefix, key)
	if err != nil {
		return errors.Wrapf(err, "build cache key err, key is %+v", key)
	}
	m.Store.Set(cacheKey, buf, int64(expiration))
	return nil
}

// Get data
func (m *memoryCache) Get(key string, val interface{}) error {
	cacheKey, err := BuildCacheKey(m.KeyPrefix, key)
	if err != nil {
		return errors.Wrapf(err, "build cache key err, key is %+v", key)
	}
	data, ok := m.Store.Get(cacheKey)
	if !ok {
		return nil
	}
	err = Unmarshal(m.encoding, data.([]byte), val)
	if err != nil {
		return errors.Wrapf(err, "unmarshal data error, key=%s, cacheKey=%s type=%v, json is %+v ",
			key, cacheKey, reflect.TypeOf(val), string(data.([]byte)))
	}
	return nil
}

// Del 删除
func (m *memoryCache) Del(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	key := keys[0]
	cacheKey, err := BuildCacheKey(m.KeyPrefix, key)
	if err != nil {
		log.Warnf("build cache key err: %+v, key is %+v", err, key)
		return err
	}
	m.Store.Del(cacheKey)
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

// Incr 自增
func (m *memoryCache) Incr(key string, step int64) (int64, error) {
	panic("implement me")
}

// Decr 自减
func (m *memoryCache) Decr(key string, step int64) (int64, error) {
	panic("implement me")
}
