package cache

import (
	"reflect"
	"time"

	"github.com/pkg/errors"

	"github.com/dgraph-io/ristretto"
)

type ristrettoCache struct {
	Store     *ristretto.Cache
	KeyPrefix string
	encoding  Encoding
}

// NewRistrettoCache create a cache client
func NewRistrettoCache(keyPrefix string, encoding Encoding) Driver {
	config := &ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	}
	store, _ := ristretto.NewCache(config)
	return &ristrettoCache{
		Store:     store,
		KeyPrefix: keyPrefix,
		encoding:  encoding,
	}
}

// Set add
func (r ristrettoCache) Set(key string, val interface{}, expiration time.Duration) error {
	buf, err := Marshal(r.encoding, val)
	if err != nil {
		return errors.Wrapf(err, "marshal data err, value is %+v", val)
	}
	cacheKey, err := BuildCacheKey(r.KeyPrefix, key)
	if err != nil {
		return errors.Wrapf(err, "build cache key err, key is %+v", key)
	}
	r.Store.Set(cacheKey, buf, int64(expiration))
	return nil
}

// Get get data by key
func (r ristrettoCache) Get(key string, val interface{}) error {
	cacheKey, err := BuildCacheKey(r.KeyPrefix, key)
	if err != nil {
		return errors.Wrapf(err, "build cache key err, key is %+v", key)
	}
	data, ok := getValue(r.Store.Get(cacheKey))
	if !ok {
		return nil
	}
	err = Unmarshal(r.encoding, data, val)
	if err != nil {
		return errors.Wrapf(err, "unmarshal data error, key=%s, cacheKey=%s type=%v, json is %+v ",
			key, cacheKey, reflect.TypeOf(val), string(data))
	}
	return nil
}

func (r ristrettoCache) MultiSet(valMap map[string]interface{}, expiration time.Duration) error {
	panic("implement me")
}

func (r ristrettoCache) MultiGet(keys []string, valueMap interface{}) error {
	panic("implement me")
}

func (r ristrettoCache) Del(keys ...string) error {
	panic("implement me")
}

func (r ristrettoCache) Incr(key string, step int64) (int64, error) {
	panic("implement me")
}

func (r ristrettoCache) Decr(key string, step int64) (int64, error) {
	panic("implement me")
}
