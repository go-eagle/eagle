package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-eagle/eagle/pkg/encoding"
)

func Test_NewMemoryCache(t *testing.T) {
	asserts := assert.New(t)

	client := NewMemoryCache("memory-unit-test", encoding.JSONEncoding{})
	asserts.NotNil(client)
}

func TestMemoStore_Set(t *testing.T) {
	asserts := assert.New(t)

	store := NewMemoryCache("memory-unit-test", encoding.JSONEncoding{})
	err := store.Set("test-key", "test-val", -1)
	asserts.NoError(err)
}

func TestMemoStore_Get(t *testing.T) {
	asserts := assert.New(t)
	store := NewMemoryCache("memory-unit-test", encoding.JSONEncoding{})

	// 正常情况
	{
		var gotVal string
		setVal := "test-val"
		err := store.Set("test-get-key", setVal, 3600)
		asserts.NoError(err)
		err = store.Get("test-get-key", &gotVal)
		asserts.NoError(err)
		t.Log(setVal, gotVal)
		asserts.Equal(setVal, gotVal)
	}
}
