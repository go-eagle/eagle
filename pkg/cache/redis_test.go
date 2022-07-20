package cache

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-eagle/eagle/pkg/encoding"
	redis2 "github.com/go-eagle/eagle/pkg/redis"
)

func Test_redisCache_SetGet(t *testing.T) {
	// 实例化redis客户端
	redis2.InitTestRedis()

	// 获取redis客户端
	redisClient := redis2.RedisClient
	// 实例化redis cache
	cache := NewRedisCache(redisClient, "unit-test", encoding.JSONEncoding{}, func() interface{} {
		return new(int64)
	})

	ctx := context.Background()

	// test set
	type setArgs struct {
		key        string
		value      interface{}
		expiration time.Duration
	}

	value := "val-001"
	setTests := []struct {
		name    string
		cache   Cache
		args    setArgs
		wantErr bool
	}{
		{
			"test redis set",
			cache,
			setArgs{"key-001", &value, 60 * time.Second},
			false,
		},
	}

	for _, tt := range setTests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.cache
			if err := c.Set(ctx, tt.args.key, tt.args.value, tt.args.expiration); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// test get
	type args struct {
		key string
	}

	tests := []struct {
		name    string
		cache   Cache
		args    args
		wantVal interface{}
		wantErr bool
	}{
		{
			"test redis get",
			cache,
			args{"key-001"},
			"val-001",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.cache
			var gotVal interface{}
			err := c.Get(ctx, tt.args.key, &gotVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log("gotval", gotVal)
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Get() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}
