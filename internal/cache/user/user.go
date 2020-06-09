package user

import (
	"fmt"
	"time"

	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/pkg/cache"
)

const (
	// PrefixUserCacheKey cache前缀
	PrefixUserCacheKey = "user:cache:%d"
	// DefaultExpireTime 默认过期时间
	DefaultExpireTime = time.Hour * 24
)

// Cache cache
type Cache struct {
	cache cache.Driver
}

// NewUserCache new一个用户cache
func NewUserCache() *Cache {
	// todo: cache.Init() 已经在main.go执行，这里应该不需要再初始化，待排查
	cache.Init()
	return &Cache{
		cache: cache.Client,
	}
}

// SetUserCache 写入用户cache
func (u *Cache) SetUserCache(userID uint64, user *model.UserModel) error {
	if user == nil || user.ID == 0 {
		return nil
	}
	cacheKey := fmt.Sprintf(PrefixUserCacheKey, userID)
	err := u.cache.Set(cacheKey, user, DefaultExpireTime)
	if err != nil {
		return err
	}
	return nil
}

// GetUserCache 获取用户cache
func (u *Cache) GetUserCache(userID uint64) (userModel *model.UserModel, err error) {
	cacheKey := fmt.Sprintf(PrefixUserCacheKey, userID)
	err = u.cache.Get(cacheKey, &userModel)
	if err != nil {
		return userModel, err
	}
	return userModel, nil
}

// MultiGetUserCache 批量获取用户cache
func (u *Cache) MultiGetUserCache(userIDs []uint64) (userMap map[uint64]*model.UserModel, err error) {
	var keys []string
	for _, v := range userIDs {
		cacheKey := fmt.Sprintf(PrefixUserCacheKey, v)
		keys = append(keys, cacheKey)
	}
	err = u.cache.MultiGet(keys, userMap)
	if err != nil {
		return nil, err
	}
	return userMap, nil
}

// DelUserCache 删除用户cache
func (u *Cache) DelUserCache(userID uint64) error {
	cacheKey := fmt.Sprintf(PrefixUserCacheKey, userID)
	err := u.cache.Del(cacheKey)
	if err != nil {
		return err
	}
	return nil
}
