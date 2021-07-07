package cache

import (
	"context"

	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/pkg/cache"
	"github.com/1024casts/snake/pkg/encoding"
	"github.com/1024casts/snake/pkg/redis"
)

func getCacheClient(ctx context.Context) cache.Cache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	client := cache.NewRedisCache(redis.RedisClient, cachePrefix, jsonEncoding, func() interface{} {
		return &model.UserBaseModel{}
	})

	return client
}
