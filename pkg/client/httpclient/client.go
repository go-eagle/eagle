// http 客户端

package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-eagle/eagle/pkg/utils"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// see: https://github.com/iiinsomnia/gochat/blob/master/utils/http.go

const (
	// ContentTypeJSON json format
	ContentTypeJSON = "application/json; charset=utf-8"
	// ContentTypeForm form format
	ContentTypeForm = "application/x-www-form-urlencoded; charset=utf-8"

	// DefaultTimeout max exec time for a request
	DefaultTimeout = 3 * time.Second
)

// ------------------ JSON ------------------

// GetJSON get json data by get method
func GetJSON(ctx context.Context, url string, options ...Option) ([]byte, error) {
	return withJSONBody(ctx, http.MethodGet, url, nil, options...)
}

// PostJSON send json data by post method
func PostJSON(ctx context.Context, url string, data json.RawMessage, options ...Option) ([]byte, error) {
	return withJSONBody(ctx, http.MethodPost, url, data, options...)
}

func withJSONBody(ctx context.Context, method, url string, raw json.RawMessage, options ...Option) (body []byte, err error) {
	opt := defaultOptions()
	for _, o := range options {
		o(opt)
	}
	opt.header["Content-Type"] = []string{ContentTypeJSON}
	return doRequest(ctx, method, url, raw, opt)
}

// ------------------ request form ------------------

// PostForm send form data by post method
func PostForm(ctx context.Context, url string, form url.Values, options ...Option) ([]byte, error) {
	return withFormBody(ctx, http.MethodPost, url, form, options...)
}

func withFormBody(ctx context.Context, method, url string, form url.Values, options ...Option) (body []byte, err error) {
	opt := defaultOptions()
	for _, o := range options {
		o(opt)
	}
	opt.header["Content-Type"] = []string{ContentTypeForm}
	formValue := form.Encode()
	return doRequest(ctx, method, url, utils.StringToBytes(formValue), opt)
}

func doRequest(ctx context.Context, method, url string, payload []byte, opt *options) (ret []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.Wrapf(err, "[httpClient] get req err")
	}

	client := &http.Client{
		// add header and set response status for tracing
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   opt.timeout,
	}

	// set header
	for key, value := range opt.header {
		req.Header.Set(key, value[0])
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "[httpClient] do request from [%s %s] err", method, url)
	}
	if resp != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}

	if !isSuccess(resp.StatusCode) {
		return nil, errors.Errorf("[httpClient] status code is %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "[httpClient] read resp body from [%s %s] err", method, url)
	}
	return body, nil
}

// isSuccess check is success
func isSuccess(statusCode int) bool {
	return statusCode < http.StatusBadRequest
}
