package http

import (
	"time"

	"github.com/go-resty/resty/v2"
)

var headerDataType = "application/json"

type restyClient struct {
}

func newRestyClient() *restyClient {
	return &restyClient{}
}

func (r restyClient) Get(url string, duration time.Duration) ([]byte, error) {
	client := resty.New()

	resp, err := client.R().
		// here can add set header
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
