package cache

import (
	"context"

	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/pkg/cache"
	"github.com/go-eagle/eagle/pkg/encoding"
	"github.com/go-eagle/eagle/pkg/redis"
)

func getCacheClient(ctx context.Context) cache.Cache {
	sonicEncoding := encoding.SonicEncoding{}
	cachePrefix := ""
	client := cache.NewRedisCache(redis.RedisClient, cachePrefix, sonicEncoding, func() interface{} {
		return &model.UserBaseModel{}
	})

	return client
}
