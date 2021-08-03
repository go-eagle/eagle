package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

// raw 使用原生包封装的 http client

// rawClient
type rawClient struct{}

// newRawClient 实例化 http 客户端
func newRawClient() Client {
	return &rawClient{}
}

// Get get data by get method
func (r *rawClient) Get(ctx context.Context, url string, params map[string]string, duration time.Duration, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   duration,
	}

	tr := otel.GetTracerProvider().Tracer("tracer from http client")
	ctx, span := tr.Start(ctx, "http request")
	defer span.End()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

// Post send data by post method
func (r *rawClient) Post(ctx context.Context, url string, data []byte, duration time.Duration, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   duration,
	}
	req.Header.Set("Content-Type", contentTypeJSON)

	tr := otel.GetTracerProvider().Tracer("tracer from http client")
	ctx, span := tr.Start(ctx, "http request")
	defer span.End()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
