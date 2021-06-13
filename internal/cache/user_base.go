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
	// PrefixUserBaseCacheKey cache前缀, 规则：业务+模块+{ID}
	PrefixUserBaseCacheKey = "snake:user:base:%d"
)

// Cache cache
type Cache struct {
	cache cache.Driver
	//localCache cache.Driver
}

// NewUserCache new一个用户cache
func NewUserCache() *Cache {
	encoding := cache.JSONEncoding{}
	cachePrefix := ""
	return &Cache{
		cache: cache.NewRedisCache(redis.RedisClient, cachePrefix, encoding, func() interface{} {
			return &model.UserBaseModel{}
		}),
	}
}

// GetUserBaseCacheKey 获取cache key
func (c *Cache) GetUserBaseCacheKey(userID uint64) string {
	return fmt.Sprintf(PrefixUserBaseCacheKey, userID)
}

// SetUserBaseCache 写入用户cache
func (c *Cache) SetUserBaseCache(ctx context.Context, userID uint64, user *model.UserBaseModel, duration time.Duration) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cache.SetUserBaseCache")
	defer span.Finish()

	if user == nil || user.ID == 0 {
		return nil
	}
	cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, userID)
	err := c.cache.Set(cacheKey, user, duration)
	if err != nil {
		return err
	}
	return nil
}

// GetUserBaseCache 获取用户cache
func (c *Cache) GetUserBaseCache(ctx context.Context, userID uint64) (data *model.UserBaseModel, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cache.GetUserBaseCache")
	defer span.Finish()
	client := getCacheClient(ctx)

	cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, userID)
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

// MultiGetUserBaseCache 批量获取用户cache
func (c *Cache) MultiGetUserBaseCache(ctx context.Context, userIDs []uint64) (map[string]*model.UserBaseModel, error) {
	var keys []string
	for _, v := range userIDs {
		cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, v)
		keys = append(keys, cacheKey)
	}

	// 需要在这里make实例化，如果在返回参数里直接定义会报 nil map
	userMap := make(map[string]*model.UserBaseModel)
	err := c.cache.MultiGet(keys, userMap)
	if err != nil {
		return nil, err
	}
	return userMap, nil
}

// DelUserBaseCache 删除用户cache
func (c *Cache) DelUserBaseCache(ctx context.Context, userID uint64) error {
	cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, userID)
	err := c.cache.Del(cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// DelUserBaseCache 删除用户cache
func (c *Cache) SetCacheWithNotFound(ctx context.Context, userID uint64) error {
	cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, userID)
	err := c.cache.SetCacheWithNotFound(cacheKey)
	if err != nil {
		return err
	}
	return nil
}
