package cache

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

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/pkg/cache"
	"github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/pkg/redis"
)

const (
	// Prefix{{.Name}}CacheKey cache前缀, 规则：业务+模块+{ID}
	Prefix{{.Name}}CacheKey = "snake:{{.Name}}:%d"
)

// Cache cache
type Cache struct {
	cache cache.Driver
	//localCache cache.Driver
}

// New{{.Name}} new一个用户cache
func New{{.Name}}() *Cache {
	encoding := cache.JSONEncoding{}
	cachePrefix := ""
	return &Cache{
		cache: cache.NewRedisCache(redis.RedisClient, cachePrefix, encoding, func() interface{} {
			return &model.{{.Name}}Model{}
		}),
	}
}

// Get{{.Name}}CacheKey 获取cache key
func (c *Cache) Get{{.Name}}CacheKey(userID uint64) string {
	return fmt.Sprintf(Prefix{{.Name}}CacheKey, userID)
}

// Set{{.Name}}Cache 写入用户cache
func (c *Cache) Set{{.Name}}Cache(ctx context.Context, userID uint64, user *model.{{.Name}}Model, duration time.Duration) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cache.Set{{.Name}}Cache")
	defer span.Finish()

	if user == nil || user.ID == 0 {
		return nil
	}
	cacheKey := fmt.Sprintf(Prefix{{.Name}}CacheKey, userID)
	err := c.cache.Set(cacheKey, user, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get{{.Name}}Cache 获取用户cache
func (c *Cache) Get{{.Name}}Cache(ctx context.Context, userID uint64) (data *model.{{.Name}}Model, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cache.Get{{.Name}}Cache")
	defer span.Finish()
	client := getCacheClient(ctx)

	cacheKey := fmt.Sprintf(Prefix{{.Name}}CacheKey, userID)
	//err = u.cache.Get(cacheKey, &data)
	err = client.Get(cacheKey, &data)
	if err != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			ext.Error.Set(span, true)
		}
		log.WithContext(ctx).Warnf("get err from redis, err: %+v", err)
		return nil, err
	}
	return data, nil
}

// MultiGet{{.Name}}Cache 批量获取用户cache
func (c *Cache) MultiGet{{.Name}}Cache(ctx context.Context, userIDs []uint64) (map[string]*model.{{.Name}}Model, error) {
	var keys []string
	for _, v := range userIDs {
		cacheKey := fmt.Sprintf(Prefix{{.Name}}CacheKey, v)
		keys = append(keys, cacheKey)
	}

	// 需要在这里make实例化，如果在返回参数里直接定义会报 nil map
	userMap := make(map[string]*model.{{.Name}}Model)
	err := c.cache.MultiGet(keys, userMap)
	if err != nil {
		return nil, err
	}
	return userMap, nil
}

// Del{{.Name}}Cache 删除用户cache
func (c *Cache) Del{{.Name}}Cache(ctx context.Context, userID uint64) error {
	cacheKey := fmt.Sprintf(Prefix{{.Name}}CacheKey, userID)
	err := c.cache.Del(cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// Del{{.Name}}Cache 删除用户cache
func (c *Cache) SetCacheWithNotFound(ctx context.Context, userID uint64) error {
	cacheKey := fmt.Sprintf(Prefix{{.Name}}CacheKey, userID)
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
