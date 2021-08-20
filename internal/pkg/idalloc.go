// Package pkg ID 分配器，主要使用redis进行分配
package pkg

import "github.com/go-eagle/eagle/pkg/redis"

// IDAlloc define struct
type IDAlloc struct {
	idGenerator *redis.IDAlloc
}

// NewIDAlloc create a id alloc
func NewIDAlloc() *IDAlloc {
	return &IDAlloc{
		idGenerator: redis.NewIDAlloc(redis.RedisClient),
	}
}

// GetUserID generate user id from redis
func (i *IDAlloc) GetUserID() (int64, error) {
	return i.idGenerator.GetNewID("user_id", 1)
}
