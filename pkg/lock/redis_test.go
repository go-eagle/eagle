package lock_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-eagle/eagle/pkg/lock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisLock(t *testing.T) {
	// Set up a miniredis server
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer s.Close()

	// Create a Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	ctx := context.Background()

	// Test acquiring the lock
	l := lock.NewRedisLock(rdb, "test-key", 5*time.Second)
	locked, err := l.Lock(ctx)
	assert.NoError(t, err, "expected no error on acquiring lock")
	assert.True(t, locked, "expected to successfully acquire the lock")

	// Test re-acquiring the lock
	lockedAgain, err := l.Lock(ctx)
	assert.NoError(t, err, "expected no error on re-acquiring lock")
	assert.True(t, lockedAgain, "expected not to acquire the lock again")

	// Test releasing the lock
	unlocked, err := l.Unlock(ctx)
	assert.NoError(t, err, "expected no error on releasing lock")
	assert.True(t, unlocked, "expected to successfully release the lock")

	// Test re-acquiring the lock after release
	lockedAgain, err = l.Lock(ctx)
	assert.NoError(t, err, "expected no error on re-acquiring lock after release")
	assert.True(t, lockedAgain, "expected to successfully re-acquire the lock after release")
}

func TestAutoRenew(t *testing.T) {
	// Set up a miniredis server
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer s.Close()

	// Create a Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Test auto-renewal of the lock
	l := lock.NewRedisLock(rdb, "test-key-auto-renew", 1*time.Second)
	locked, err := l.Lock(ctx)
	assert.NoError(t, err, "expected no error on acquiring lock")
	assert.True(t, locked, "expected to successfully acquire the lock")

	// Wait for some time to ensure the lock is being renewed
	time.Sleep(1500 * time.Millisecond)

	// Check the TTL of the key to verify that it has been renewed
	ttl, err := l.GetTTL(ctx)
	assert.NoError(t, err, "expected no error on checking TTL")
	assert.Greater(t, ttl, time.Duration(0), "expected TTL to be greater than zero due to renewal")

	// Release the lock
	unlocked, err := l.Unlock(ctx)
	assert.NoError(t, err, "expected no error on releasing lock")
	assert.True(t, unlocked, "expected to successfully release the lock")
}
