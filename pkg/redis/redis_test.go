package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	InitTestRedis()

	// test init
	ret, err := RedisClient.Ping(context.Background()).Result()
	assert.Nil(t, err)
	assert.Equal(t, "PONG", ret)
}

func TestNewRedisClient(t *testing.T) {
	// test default client
	InitTestRedis()
	client := NewRedisManager()
	assert.NotNil(t, client)

	// get a not exist client
	// rdb, err := client.GetClient("not-exist")
	// assert.NotNil(t, err)
	// assert.Nil(t, rdb)
}

func TestRedisSetGet(t *testing.T) {
	InitTestRedis()

	var setGetKey = "test-set"
	var setGetValue = "test-content"
	RedisClient.Set(context.Background(), setGetKey, setGetValue, time.Second*100)

	expectValue := RedisClient.Get(context.Background(), setGetKey).Val()
	assert.Equal(t, setGetValue, expectValue)
}
