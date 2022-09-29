package lock

import (
	"context"
)

const (
	// RedisLockKey redis lock key
	RedisLockKey = "eagle:redis:lock:%s"
	// EtcdLockKey etcd lock key
	EtcdLockKey = "/eagle/lock/%s"
)

// Lock define common func
type Lock interface {
	Lock(ctx context.Context) (bool, error)
	Unlock(ctx context.Context) (bool, error)
}
