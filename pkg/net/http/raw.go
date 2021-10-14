package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// raw 使用原生包封装的 http client

// rawClient
type rawClient struct {
	resp *http.Response
}

// newRawClient 实例化 http 客户端
func newRawClient() *rawClient {
	return &rawClient{}
}

// Get get data by get method
func (r *rawClient) GetJSON(ctx context.Context, url string, options ...Option) ([]byte, error) {
	return r.withJSONBody(ctx, http.MethodGet, url, nil, options...)
}

// Post send data by post method
func (r *rawClient) PostJSON(ctx context.Context, url string, data json.RawMessage, options ...Option) ([]byte, error) {
	return r.withJSONBody(ctx, http.MethodPost, url, data, options...)
}

func (r *rawClient) withJSONBody(ctx context.Context, method, url string, raw json.RawMessage, options ...Option) (body []byte, err error) {
	opt := &option{}
	for _, f := range options {
		f(opt)
	}
	opt.header["Content-Type"] = []string{"application/json; charset=utf-8"}
	return doRequest(ctx, method, url, raw, opt)
}

func doRequest(ctx context.Context, method, url string, payload []byte, opt *option) (ret []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.Wrapf(err, "[httpClient] get req err")
	}

	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   opt.timeout,
	}

	// set header
	for key, value := range opt.header {
		req.Header.Set(key, value[0])
	}

	ctx, span := tracer.Start(req.Context(), "HTTP Post")
	defer span.End()

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "[httpClient] do request from [%s %s] err", method, url)
	}
	if resp != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}

	if !isSuccess(resp) {
		return nil, errors.New("request")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "[httpClient] read resp body from [%s %s] err", method, url)
	}
	return body, nil
}

// isSuccess check is success
func isSuccess(resp *http.Response) bool {
	return resp.StatusCode < http.StatusBadRequest
}
