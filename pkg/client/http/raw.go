package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/glog"
)

const (
	DefaultTimeOut = 1 * time.Second
)

// raw 使用原生包封装的 http client

type rawClient struct{}

func newRawClient() Client {
	return &rawClient{}
}

func (r *rawClient) Get(url string, params map[string]string, duration time.Duration) ([]byte, error) {
	client := http.Client{Timeout: DefaultTimeOut}
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Error(err)
		return err
	}

	if err := json.Unmarshal(b, target); err != nil {
		glog.Error(err)
		return fmt.Errorf("can't unmarshal to target %s %s", err, b)
	}

	return nil
}

func (r *rawClient) Post(url string, requestBody string, duration time.Duration) ([]byte, error) {
	panic("implement me")
}

func (r *rawClient) PostJson(url string, requestBody string, duration time.Duration) ([]byte, error) {
	panic("implement me")
}
