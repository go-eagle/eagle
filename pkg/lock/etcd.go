package lock

import (
	"context"
	"fmt"
	"time"

	v3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// EtcdLock define a etcd lock
type EtcdLock struct {
	sess *concurrency.Session
	mu   *concurrency.Mutex
}

// NewEtcdLock create a etcd lock
// ttl for lease
func NewEtcdLock(client *v3.Client, key string, ttl int) (mutex *EtcdLock, err error) {
	mutex = &EtcdLock{}
	// default session ttl = 60s
	mutex.sess, err = concurrency.NewSession(client, concurrency.WithTTL(ttl))
	if err != nil {
		return
	}
	mutex.mu = concurrency.NewMutex(mutex.sess, getEtcdKey(key))
	return
}

// Lock acquires the lock.
func (l *EtcdLock) Lock(ctx context.Context, timeout time.Duration) (b bool, err error) {
	// get lock timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	// NOTE: ignore bool value
	return true, l.mu.Lock(ctx)
}

// Unlock release a lock.
func (l *EtcdLock) Unlock(ctx context.Context) (b bool, err error) {
	err = l.mu.Unlock(ctx)
	if err != nil {
		return
	}
	// NOTE: ignore bool value
	return true, l.sess.Close()
}

// getEtcdKey 获取key
func getEtcdKey(key string) string {
	return fmt.Sprintf(EtcdLockKey, key)
}
