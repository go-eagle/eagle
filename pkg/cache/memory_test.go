package cache

import (
	"reflect"
	"testing"
	"time"
)

func Test_memoryCache_SetGet(t *testing.T) {
	// 实例化memory cache
	cache := NewMemoryCache("memory-unit-test", JSONEncoding{})

	// test set
	type setArgs struct {
		key        string
		value      interface{}
		expiration time.Duration
	}

	setTests := []struct {
		name    string
		cache   Driver
		args    setArgs
		wantErr bool
	}{
		{
			"test memory set",
			cache,
			setArgs{"key-001", "{\"username\":\"snake\"}", 60 * time.Second},
			false,
		},
	}

	for _, tt := range setTests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.cache
			if err := c.Set(tt.args.key, tt.args.value, tt.args.expiration); (err != nil) != tt.wantErr {
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
		cache   Driver
		args    args
		wantVal interface{}
		wantErr bool
	}{
		{
			"test memory get",
			cache,
			args{"key-001"},
			"{\"username\":\"snake\"}",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.cache
			gotVal, err := c.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Get() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}
