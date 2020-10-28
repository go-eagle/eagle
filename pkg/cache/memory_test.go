package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewMemoryCache(t *testing.T) {
	asserts := assert.New(t)

	client := NewMemoryCache("memory-unit-test", JSONEncoding{})
	asserts.NotNil(client)
}

func TestMemoStore_Set(t *testing.T) {
	asserts := assert.New(t)

	store := NewMemoryCache("unit-test", JSONEncoding{})
	err := store.Set("test-key", "test-val", -1)
	asserts.NoError(err)

	val, ok := store.Store.Load("unit-test:test-key")
	asserts.True(ok)
	asserts.Equal("test-val", val.(itemWithTTL).value)

	store.GarbageCollect()
}

func TestMemoStore_Get(t *testing.T) {
	asserts := assert.New(t)
	store := NewMemoryCache("unit-test", JSONEncoding{})

	// 正常情况
	{
		var val string
		err := store.Set("test-key", "test-val", -1)
		asserts.NoError(err)
		err = store.Get("test-key", &val)
		t.Log("val.......", val)
		asserts.NoError(err)
		asserts.Equal("test-val", val)
	}

	// Key不存在
	{
		var val string
		err := store.Get("something", val)
		asserts.Error(err)
	}

	// 过期
	{
		var val string
		err := store.Set("test-key", "test-val", 1)
		asserts.NoError(err)
		err = store.Get("test-key", val)
		asserts.NotEmpty(val)
	}

}
