package add

import (
	"bytes"
	"strings"

	"github.com/alecthomas/template"
)

const cacheTemplate = `
package cache

//go:generate mockgen -source=internal/dal/cache/{{.UsName}}_cache.go -destination=internal/mock/{{.UsName}}_cache_mock.go  -package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/go-eagle/eagle/pkg/cache"
	"github.com/go-eagle/eagle/pkg/encoding"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/utils"
	"github.com/redis/go-redis/v9"

	"{{.ModName}}/internal/dal/db/model"
)

var (
	// Prefix{{.Name}}CacheKey cache prefix
	Prefix{{.Name}}CacheKey = utils.ConcatString(prefix, "{{.ColonName}}:%d")
)

// {{.Name}}Cache define cache interface
type {{.Name}}Cache interface {
	Set{{.Name}}Cache(ctx context.Context, id int64, data *model.{{.Name}}Model, duration time.Duration) error
	Get{{.Name}}Cache(ctx context.Context, id int64) (data *model.{{.Name}}Model, err error)
	MultiGet{{.Name}}Cache(ctx context.Context, ids []int64) (map[string]*model.{{.Name}}Model, error)
	MultiSet{{.Name}}Cache(ctx context.Context, data []*model.{{.Name}}Model, duration time.Duration) error
	Del{{.Name}}Cache(ctx context.Context, id int64) error
	SetCacheWithNotFound(ctx context.Context, id int64) error
}

// {{.LcName}}Cache define cache struct
type {{.LcName}}Cache struct {
	cache cache.Cache
}

// New{{.Name}}Cache new a cache
func New{{.Name}}Cache(rdb *redis.Client) {{.Name}}Cache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	return &{{.LcName}}Cache{
		cache: cache.NewRedisCache(rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.{{.Name}}Model{}
		}),
	}
}

// Get{{.Name}}CacheKey get cache key
func (c *{{.LcName}}Cache) Get{{.Name}}CacheKey(id int64) string {
	return fmt.Sprintf(Prefix{{.Name}}CacheKey, id)
}

// Set{{.Name}}Cache write to cache
func (c *{{.LcName}}Cache) Set{{.Name}}Cache(ctx context.Context, id int64, data *model.{{.Name}}Model, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.Get{{.Name}}CacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get{{.Name}}Cache get from cache
func (c *{{.LcName}}Cache) Get{{.Name}}Cache(ctx context.Context, id int64) (data *model.{{.Name}}Model, err error) {
	cacheKey := c.Get{{.Name}}CacheKey(id)
	err = c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		log.WithContext(ctx).Warnf("[cache] Get{{.Name}}Cache err from redis, err: %+v", err)
		return nil, err
	}
	return data, nil
}

// MultiGet{{.Name}}Cache batch get cache
func (c *{{.LcName}}Cache) MultiGet{{.Name}}Cache(ctx context.Context, ids []int64) (map[string]*model.{{.Name}}Model, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.Get{{.Name}}CacheKey(v)
		keys = append(keys, cacheKey)
	}

	// NOTE: 需要在这里make实例化，如果在返回参数里直接定义会报 nil map
	retMap := make(map[string]*model.{{.Name}}Model)
	err := c.cache.MultiGet(ctx, keys, retMap)
	if err != nil {
		return nil, err
	}
	return retMap, nil
}

// MultiSet{{.Name}}Cache batch set cache
func (c *{{.LcName}}Cache) MultiSet{{.Name}}Cache(ctx context.Context, data []*model.{{.Name}}Model, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.Get{{.Name}}CacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}
	return nil
}

// Del{{.Name}}Cache delete cache
func (c *{{.LcName}}Cache) Del{{.Name}}Cache(ctx context.Context, id int64) error {
	cacheKey := c.Get{{.Name}}CacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *{{.LcName}}Cache) SetCacheWithNotFound(ctx context.Context, id int64) error {
	cacheKey := c.Get{{.Name}}CacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
`

func (c *Cache) execute() ([]byte, error) {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("cache").Parse(strings.TrimSpace(cacheTemplate))
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, c); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
