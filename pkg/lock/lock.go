package lock

import (
	"context"
	"time"
)

const (
	// LockKey redis lock key
	LockKey = "eagle:redis:lock:%s"
	// DefaultTimeout default expire time
	DefaultTimeout = 2 * time.Second
)

type Lock interface {
	Lock(ctx context.Context) (bool, error)
	Unlock(ctx context.Context) error
}
