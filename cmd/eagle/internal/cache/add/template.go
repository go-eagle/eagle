package add

import (
	"bytes"
	"strings"

	"github.com/alecthomas/template"
)

const cacheTemplate = `
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-eagle/eagle/pkg/cache"
	"github.com/go-eagle/eagle/pkg/encoding"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/redis"

	"{{.ModName}}/internal/model"
)

const (
	// Prefix{{.Name}}CacheKey cache prefix
	Prefix{{.Name}}CacheKey = "{{.Name}}:%d"
)

// {{.Name}}Cache define a cache struct
type {{.Name}}Cache struct {
	cache cache.Cache
}

// New{{.Name}} new a cache
func New{{.Name}}() *Cache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	return &Cache{
		cache: cache.NewRedisCache(redis.RedisClient, cachePrefix, jsonEncoding, func() interface{} {
			return &model.{{.Name}}Model{}
		}),
	}
}

// Get{{.Name}}CacheKey get cache key
func (c *{{.Name}}Cache) Get{{.Name}}CacheKey(id int64) string {
	return fmt.Sprintf(Prefix{{.Name}}CacheKey, id)
}

// Set{{.Name}}Cache write to cache
func (c *{{.Name}}Cache) Set{{.Name}}Cache(ctx context.Context, id int64, data *model.{{.Name}}Model, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.Get{{.Name}}CacheKey(id)
	err := c.cache.Set(cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get{{.Name}}Cache 获取cache
func (c *{{.Name}}Cache) Get{{.Name}}Cache(ctx context.Context, id int64) (data *model.{{.Name}}Model, err error) {
	cacheKey := c.Get{{.Name}}CacheKey(id)
	err = c.cache.Get(cacheKey, &data)
	if err != nil {
		log.WithContext(ctx).Warnf("get err from redis, err: %+v", err)
		return nil, err
	}
	return data, nil
}

// MultiGet{{.Name}}Cache 批量获取cache
func (c *{{.Name}}Cache) MultiGet{{.Name}}Cache(ctx context.Context, ids []int64) (map[string]*model.{{.Name}}Model, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.Get{{.Name}}CacheKey(v)
		keys = append(keys, cacheKey)
	}

	// NOTE: 需要在这里make实例化，如果在返回参数里直接定义会报 nil map
	retMap := make(map[string]*model.{{.Name}}Model)
	err := c.cache.MultiGet(keys, retMap)
	if err != nil {
		return nil, err
	}
	return retMap, nil
}

// Del{{.Name}}Cache 删除cache
func (c *{{.Name}}Cache) Del{{.Name}}Cache(ctx context.Context, id int64) error {
	cacheKey := c.Get{{.Name}}CacheKey(id)
	err := c.cache.Del(cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// Del{{.Name}}Cache set empty cache
func (c *{{.Name}}Cache) SetCacheWithNotFound(ctx context.Context, id int64) error {
	cacheKey := c.Get{{.Name}}CacheKey(id)
	err := c.cache.SetCacheWithNotFound(cacheKey)
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
