package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/1024casts/snake/pkg/log"
)

// raw 使用原生包封装的 http client

// rawClient
type rawClient struct{}

// newRawClient 实例化 http 客户端
func newRawClient() Client {
	return &rawClient{}
}

// Get get data by get method
func (r *rawClient) Get(url string, params map[string]string, duration time.Duration) ([]byte, error) {
	client := http.Client{Timeout: duration}
	var target []byte

	resp, err := client.Get(url)
	if err != nil {
		log.Warnf("get url:%s, err: %s", url, err)
		return target, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("read body: %s by ioutil, err: %s", b, err)
		return target, err
	}

	if err := json.Unmarshal(b, &target); err != nil {
		log.Warnf("can't unmarshal to target err: %s, body: %s", err, b)
		return target, fmt.Errorf("can't unmarshal to target err: %s, body: %s", err, b)
	}

	return target, nil
}

// Post send data by post method
func (r *rawClient) Post(url string, data []byte, duration time.Duration) ([]byte, error) {
	client := http.Client{Timeout: duration}
	var target []byte
	resp, err := client.Post(url, contentTypeJson, bytes.NewBuffer(data))
	if err != nil {
		log.Warnf("post url:%s, err: %s", url, err)
		return target, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("read body: %s by ioutil, err: %s", b, err)
		return target, err
	}

	log.Infof("resp: %+v", string(b))
	if err := json.Unmarshal(b, &target); err != nil {
		log.Warnf("can't unmarshal to target err: %s, body: %s", err, b)
		return target, fmt.Errorf("can't unmarshal to target, err: %s, body: %s", err, b)
	}

	return target, nil
}
