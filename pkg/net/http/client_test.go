package http

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpClient(t *testing.T) {
	t.Run("test http get json func", func(t *testing.T) {
		var ret []byte
		var want = "http://httpbin.org/get"
		ret, err := GetJSON(context.Background(), "http://httpbin.org/get")
		if err != nil {
			t.Log(err)
		}
		type resp struct {
			Url string `json:"url"`
		}
		r := resp{}

		err = json.Unmarshal(ret, &r)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, r.Url, want)
	})

	t.Run("test http post json func", func(t *testing.T) {
		var ret []byte
		jsonStr := `{"key1":"value1"}`
		ret, err := PostJSON(context.Background(), "http://httpbin.org/post", []byte(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		type resp struct {
			Data string `json:"data"`
		}
		r := resp{}

		err = json.Unmarshal(ret, &r)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, r.Data, jsonStr)
	})
}
