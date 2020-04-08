// http客户端 resty

package http

import (
	"time"

	"github.com/go-resty/resty/v2"
)

// docs: https://github.com/go-resty/resty

const (
	headerDataType = "application/json"
)

type restyClient struct {
}

func newRestyClient() Client {
	return &restyClient{}
}

func (r restyClient) Get(url string, params map[string]string, duration time.Duration) ([]byte, error) {
	client := resty.New()

	if duration != 0 {
		client.SetTimeout(duration)
	}

	if len(params) > 0 {
		client.SetPathParams(params)
	}

	resp, err := client.R().
		SetHeaders(map[string]string{
			"Content-Type": headerDataType,
		}).
		Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body(), nil
}

func (r restyClient) Post(url string, requestBody string, duration time.Duration) ([]byte, error) {
	client := resty.New()

	if duration != 0 {
		client.SetTimeout(duration)
	}

	cr := client.R().
		SetBody(requestBody).
		SetHeaders(map[string]string{
			"Content-Type": headerDataType,
		})

	resp, err := cr.Post(url)
	if err != nil {
		return nil, err
	}

	return resp.Body(), nil
}
