package user

import (
	"context"
	"fmt"
	"time"

	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/pkg/cache"
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
func (u *Cache) GetUserBaseCacheKey(userID uint64) string {
	return fmt.Sprintf(PrefixUserBaseCacheKey, userID)
}

// SetUserBaseCache 写入用户cache
func (u *Cache) SetUserBaseCache(ctx context.Context, userID uint64, user *model.UserBaseModel, duration time.Duration) error {
	if user == nil || user.ID == 0 {
		return nil
	}
	cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, userID)
	err := u.cache.Set(cacheKey, user, duration)
	if err != nil {
		return err
	}
	return nil
}

// GetUserBaseCache 获取用户cache
func (u *Cache) GetUserBaseCache(ctx context.Context, userID uint64) (data *model.UserBaseModel, err error) {
	cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, userID)
	err = u.cache.Get(cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiGetUserBaseCache 批量获取用户cache
func (u *Cache) MultiGetUserBaseCache(ctx context.Context, userIDs []uint64) (map[string]*model.UserBaseModel, error) {
	var keys []string
	for _, v := range userIDs {
		cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, v)
		keys = append(keys, cacheKey)
	}

	// 需要在这里make实例化，如果在返回参数里直接定义会报 nil map
	userMap := make(map[string]*model.UserBaseModel)
	err := u.cache.MultiGet(keys, userMap)
	if err != nil {
		return nil, err
	}
	return userMap, nil
}

// DelUserBaseCache 删除用户cache
func (u *Cache) DelUserBaseCache(ctx context.Context, userID uint64) error {
	cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, userID)
	err := u.cache.Del(cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// DelUserBaseCache 删除用户cache
func (u *Cache) SetCacheWithNotFound(ctx context.Context, userID uint64) error {
	cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, userID)
	err := u.cache.SetCacheWithNotFound(cacheKey)
	if err != nil {
		return err
	}
	return nil
}
