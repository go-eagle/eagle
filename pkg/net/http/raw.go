package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-eagle/eagle/pkg/log"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// raw 使用原生包封装的 http client

// rawClient
type rawClient struct {
	resp *http.Response
}

// newRawClient 实例化 http 客户端
func newRawClient() Client {
	return &rawClient{}
}

// Get get data by get method
func (r *rawClient) Get(ctx context.Context, url string, params map[string]string, duration time.Duration, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return err
	}

	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   duration,
	}

	ctx, span := tracer.Start(ctx, "HTTP Get")
	defer span.End()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	r.resp = resp

	if !r.isSuccess() {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil && err != io.ErrUnexpectedEOF {
			log.WithContext(ctx).Errorf("[httpClient] get url: %s, status code: %s, err: %s",
				url, resp.StatusCode, err.Error())
			return err
		}
		return errors.New(string(body))
	}
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

// Post send data by post method
func (r *rawClient) Post(ctx context.Context, url string, data []byte, duration time.Duration, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   duration,
	}
	req.Header.Set("Content-Type", contentTypeJSON)

	ctx, span := tracer.Start(req.Context(), "HTTP Post")
	defer span.End()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	r.resp = resp

	if !r.isSuccess() {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.WithContext(ctx).Errorf("[httpClient] post url: %s, status code: %s, err: %s",
				url, resp.StatusCode, err.Error())
			return err
		}
		return errors.New(string(body))
	}
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

// isSuccess check is success
func (r *rawClient) isSuccess() bool {
	return r.resp.StatusCode < http.StatusBadRequest
}
