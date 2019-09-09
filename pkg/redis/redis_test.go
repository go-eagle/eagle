package redis

import (
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	Init()

	var setGetKey = "test-set"
	var setGetValue = "test-content"
	Client.Set(setGetKey, setGetValue, time.Second*100)

	expectValue := Client.Get(setGetKey).Val()
	if setGetValue != expectValue {
		t.Fail()
		return
	}

	t.Log("redis set test pass")
}
