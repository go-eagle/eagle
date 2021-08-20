package cache

import (
	"context"

	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/pkg/cache"
	"github.com/go-eagle/eagle/pkg/encoding"
	"github.com/go-eagle/eagle/pkg/redis"
)

func getCacheClient(ctx context.Context) cache.Cache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	client := cache.NewRedisCache(redis.RedisClient, cachePrefix, jsonEncoding, func() interface{} {
		return &model.UserBaseModel{}
	})

	return client
}
