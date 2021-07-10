package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// raw 使用原生包封装的 http client

// rawClient
type rawClient struct{}

// newRawClient 实例化 http 客户端
func newRawClient() Client {
	return &rawClient{}
}

// Get get data by get method
func (r *rawClient) Get(url string, params map[string]string, duration time.Duration, out interface{}) error {
	client := http.Client{Timeout: duration}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &out); err != nil {
		return fmt.Errorf("can't unmarshal to out, err: %s, body: %s", err, b)
	}

	return nil
}

// Post send data by post method
func (r *rawClient) Post(url string, data []byte, duration time.Duration, out interface{}) error {
	client := http.Client{Timeout: duration}
	resp, err := client.Post(url, contentTypeJSON, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &out); err != nil {
		return fmt.Errorf("can't unmarshal to out, err: %s, body: %s", err, b)
	}

	return nil
}
