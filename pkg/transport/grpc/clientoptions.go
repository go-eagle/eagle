package grpc

import (
	"time"

	"google.golang.org/grpc"
)

// clientOptions define gRPc client options.
type clientOptions struct {
	endpoint     string
	timeout      time.Duration
	inters       []grpc.UnaryClientInterceptor
	dialOpts     []grpc.DialOption
	balancerName string
	enableGzip   bool
	enableMetric bool
	enableLog    bool
	// retry config
	disableRetry bool
	NumRetries   int // maximum number of retry attempts
}

// ClientOption is a gRPC client option.
type ClientOption func(o *clientOptions)

// WithEndpoint with a endpoint.
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// WithTimeout with a client timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

// WithMetric enable metric.
func WithMetric() ClientOption {
	return func(o *clientOptions) {
		o.enableMetric = true
	}
}

// WithLog enable log.
func WithLog() ClientOption {
	return func(o *clientOptions) {
		o.enableLog = true
	}
}

// WithGzip enable gzip.
func WithGzip() ClientOption {
	return func(o *clientOptions) {
		o.enableGzip = true
	}
}

// WithoutRetry disable retry.
func WithoutRetry() ClientOption {
	return func(o *clientOptions) {
		o.disableRetry = true
	}
}

// WithUnaryInterceptor with unary client interceptor.
func WithUnaryInterceptor(inter ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.inters = inter
	}
}

// WithOptions with gRPC dial option.
func WithOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.dialOpts = opts
	}
}
