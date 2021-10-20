package lock

import (
	"context"
	"time"
)

const (
	// RedisLockKey redis lock key
	RedisLockKey = "eagle:redis:lock:%s"
	// EtcdLockKey
	EtcdLockKey = "/eagle/lock/%s"
)

// Lock define common func
type Lock interface {
	Lock(ctx context.Context, timeout time.Duration) (bool, error)
	Unlock(ctx context.Context) (bool, error)
}
