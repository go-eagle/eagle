package repository

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/1024casts/snake/pkg/net/tracing"
	"github.com/opentracing/opentracing-go"
)

// Do executes an HTTP request and returns the response body as a string.
// Non-200 response codes will be returned as an error with the response body.
func Do(req *http.Request) (string, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("StatusCode: %d, Body: %s", resp.StatusCode, body)
	}
	return string(body), nil
}

// 调用其他服务时，注入span
func CallOtherMicroservice(ctx context.Context, hostPort string) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ping-send")
	defer span.Finish()

	url := fmt.Sprintf("http://%s/ping", hostPort)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if err := tracing.Inject(span, req); err != nil {
		return "", err
	}
	return Do(req)
}
