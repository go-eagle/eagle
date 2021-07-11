package cache

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"

	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/pkg/cache"
	"github.com/1024casts/snake/pkg/encoding"
	"github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/pkg/redis"
)

const (
	// PrefixUserBaseCacheKey cache前缀, 规则：业务+模块+{ID}
	PrefixUserBaseCacheKey = "snake:user:base:%d"
)

// Cache cache
type Cache struct {
	cache  cache.Cache
	tracer trace.Tracer
	//localCache cache.Cache
}

// NewUserCache new一个用户cache
func NewUserCache() *Cache {
	encoding := encoding.JSONEncoding{}
	cachePrefix := ""
	return &Cache{
		cache: cache.NewRedisCache(redis.RedisClient, cachePrefix, encoding, func() interface{} {
			return &model.UserBaseModel{}
		},
		),
		tracer: otel.GetTracerProvider().Tracer("cache"),
	}
}

// GetUserBaseCacheKey 获取cache key
func (c *Cache) GetUserBaseCacheKey(userID uint64) string {
	return fmt.Sprintf(PrefixUserBaseCacheKey, userID)
}

// SetUserBaseCache 写入用户cache
func (c *Cache) SetUserBaseCache(ctx context.Context, userID uint64, user *model.UserBaseModel, duration time.Duration) error {
	ctx, span := c.tracer.Start(ctx, "SetUserBaseCache")
	defer span.End()

	if user == nil || user.ID == 0 {
		return nil
	}
	cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, userID)
	err := c.cache.Set(ctx, cacheKey, user, duration)
	if err != nil {
		return err
	}
	return nil
}

// GetUserBaseCache 获取用户cache
func (c *Cache) GetUserBaseCache(ctx context.Context, userID uint64) (data *model.UserBaseModel, err error) {
	ctx, span := c.tracer.Start(ctx, "GetUserBaseCache")
	defer span.End()

	client := getCacheClient(ctx)

	cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, userID)
	//err = u.cache.Get(cacheKey, &data)
	err = client.Get(ctx, cacheKey, &data)
	if err != nil {
		log.WithContext(ctx).Warnf("get err from redis, err: %+v", err)
		return nil, err
	}
	return data, nil
}

// MultiGetUserBaseCache 批量获取用户cache
func (c *Cache) MultiGetUserBaseCache(ctx context.Context, userIDs []uint64) (map[string]*model.UserBaseModel, error) {
	ctx, span := c.tracer.Start(ctx, "MultiGetUserBaseCache")
	defer span.End()
	var keys []string
	for _, v := range userIDs {
		cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, v)
		keys = append(keys, cacheKey)
	}

	// 需要在这里make实例化，如果在返回参数里直接定义会报 nil map
	userMap := make(map[string]*model.UserBaseModel)
	err := c.cache.MultiGet(ctx, keys, userMap)
	if err != nil {
		return nil, err
	}
	return userMap, nil
}

// DelUserBaseCache 删除用户cache
func (c *Cache) DelUserBaseCache(ctx context.Context, userID uint64) error {
	ctx, span := c.tracer.Start(ctx, "DelUserBaseCache")
	defer span.End()
	cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, userID)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound .
func (c *Cache) SetCacheWithNotFound(ctx context.Context, userID uint64) error {
	ctx, span := c.tracer.Start(ctx, "SetCacheWithNotFound")
	defer span.End()
	cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, userID)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
