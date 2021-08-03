// http客户端 resty

package http

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/go-resty/resty/v2"
)

// docs: https://github.com/go-resty/resty

type restyClient struct{}

func newRestyClient() Client {
	return &restyClient{}
}

// Get request url by get method
func (r *restyClient) Get(ctx context.Context, url string, params map[string]string, duration time.Duration, out interface{}) error {
	client := resty.New()

	if duration != 0 {
		client.SetTimeout(duration)
	}

	if params != nil {
		client.SetQueryParams(params)
	}

	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type": contentTypeJSON,
		}).
		Get(url)
	if err != nil {
		return err
	}
	defer resp.RawResponse.Body.Close()

	if resp.StatusCode() >= 400 {
		body, err := ioutil.ReadAll(resp.RawBody())
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	decoder := json.NewDecoder(resp.RawBody())
	return decoder.Decode(out)
}

// Post request url by post method
func (r *restyClient) Post(ctx context.Context, url string, data []byte, duration time.Duration, out interface{}) error {
	client := resty.New()

	if duration != 0 {
		client.SetTimeout(duration)
	}

	cr := client.R().
		SetBody(string(data)).
		SetHeaders(map[string]string{
			"Content-Type": contentTypeJSON,
		})

	resp, err := cr.Post(url)
	if err != nil {
		return err
	}
	defer resp.RawBody().Close()

	if resp.StatusCode() >= 400 {
		body, err := ioutil.ReadAll(resp.RawBody())
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	decoder := json.NewDecoder(resp.RawBody())
	return decoder.Decode(out)
}
