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
	sess       *concurrency.Session
	mu         *concurrency.Mutex
	cancelFunc context.CancelFunc
}

// NewEtcdLock create a etcd lock
// ttl for lease
func NewEtcdLock(client *v3.Client, key string, ttl int) (mutex *EtcdLock, err error) {
	mutex = &EtcdLock{}

	// get lock timeout
	// set lease ttl == request timeout
	expiration := time.Duration(ttl) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), expiration)
	mutex.cancelFunc = cancel

	// default session ttl = 60s
	mutex.sess, err = concurrency.NewSession(
		client,
		concurrency.WithTTL(ttl),
		concurrency.WithContext(ctx),
	)
	if err != nil {
		return
	}

	mutex.mu = concurrency.NewMutex(mutex.sess, getEtcdKey(key))

	return
}

// Lock acquires the lock.
func (l *EtcdLock) Lock(ctx context.Context) (b bool, err error) {
	// NOTE: ignore bool value
	return true, l.mu.Lock(ctx)
}

// Unlock release a lock.
func (l *EtcdLock) Unlock(ctx context.Context) (b bool, err error) {
	defer l.cancelFunc()

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
