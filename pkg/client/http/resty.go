// http客户端 resty

package http

import (
	"time"

	"github.com/1024casts/snake/pkg/log"

	"github.com/go-resty/resty/v2"
)

// docs: https://github.com/go-resty/resty

type restyClient struct{}

func newRestyClient() Client {
	return &restyClient{}
}

// Get request url by get method
func (r *restyClient) Get(url string, params map[string]string, duration time.Duration) ([]byte, error) {
	client := resty.New()

	if duration != 0 {
		client.SetTimeout(duration)
	}

	if len(params) > 0 {
		client.SetQueryParams(params)
	}

	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type": contentTypeJSON,
		}).
		Get(url)
	if err != nil {
		log.Warnf("get url: %s err: %s", url, err)
		return nil, err
	}

	return resp.Body(), nil
}

// Post request url by post method
func (r *restyClient) Post(url string, data []byte, duration time.Duration) ([]byte, error) {
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
		log.Warnf("post url: %s err: %s", url, err)
		return nil, err
	}

	return resp.Body(), nil
}
