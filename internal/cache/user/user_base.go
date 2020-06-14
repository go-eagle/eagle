package user

import (
	"fmt"
	"time"

	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/pkg/cache"
	"github.com/1024casts/snake/pkg/redis"
)

const (
	// PrefixUserBaseCacheKey cache前缀
	PrefixUserBaseCacheKey = "user:cache:%d"
	// DefaultExpireTime 默认过期时间
	DefaultExpireTime = time.Hour * 24
)

// Cache cache
type Cache struct {
	cache cache.Driver
}

// NewUserCache new一个用户cache
func NewUserCache() *Cache {
	encoding := cache.JSONEncoding{}
	cachePrefix := cache.PrefixCacheKey
	// TODO: redis已经全局实例化，redis.Init() 已经在main.go执行，这里应该不需要再初始化，待排查
	redis.Init()
	return &Cache{
		cache: cache.NewRedisCache(redis.Client, cachePrefix, encoding, func() interface{} {
			return &model.UserBaseModel{}
		}),
	}
}

// GetUserBaseCacheKey 获取cache key
func (u *Cache) GetUserBaseCacheKey(userID uint64) string {
	return fmt.Sprintf(cache.PrefixCacheKey+":"+PrefixUserBaseCacheKey, userID)
}

// SetUserBaseCache 写入用户cache
func (u *Cache) SetUserBaseCache(userID uint64, user *model.UserBaseModel) error {
	if user == nil || user.ID == 0 {
		return nil
	}
	cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, userID)
	err := u.cache.Set(cacheKey, user, DefaultExpireTime)
	if err != nil {
		return err
	}
	return nil
}

// GetUserBaseCache 获取用户cache
func (u *Cache) GetUserBaseCache(userID uint64) (userModel *model.UserBaseModel, err error) {
	cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, userID)
	err = u.cache.Get(cacheKey, &userModel)
	if err != nil {
		return userModel, err
	}
	return userModel, nil
}

// MultiGetUserBaseCache 批量获取用户cache
func (u *Cache) MultiGetUserBaseCache(userIDs []uint64) (map[string]*model.UserBaseModel, error) {
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
func (u *Cache) DelUserBaseCache(userID uint64) error {
	cacheKey := fmt.Sprintf(PrefixUserBaseCacheKey, userID)
	err := u.cache.Del(cacheKey)
	if err != nil {
		return err
	}
	return nil
}
