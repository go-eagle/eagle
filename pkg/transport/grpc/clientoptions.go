package grpc

import (
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/go-eagle/eagle/pkg/registry"
)

// clientOptions define gRPc client options.
type clientOptions struct {
	endpoint        string
	timeout         time.Duration
	inters          []grpc.UnaryClientInterceptor
	dialOpts        []grpc.DialOption
	balancerName    string
	discovery       registry.Discovery
	enableKeepalive bool
	kp              keepalive.ClientParameters
	enableGzip      bool
	enableMetric    bool
	enableLog       bool

	// enableTracing enables distributed tracing using OpenTelemetry protocol
	enableTracing bool
	// retry config
	disableRetry bool
	// maximum number of retry attempts
	NumRetries int
	// TracerOptions are options for OpenTelemetry gRPC interceptor.
	TracerOptions []otelgrpc.Option
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

// WithTracing enable trace.
func WithTracing() ClientOption {
	return func(o *clientOptions) {
		o.enableTracing = true
	}
}

// WithKeepalive enable keepalive.
func WithKeepalive() ClientOption {
	return func(o *clientOptions) {
		o.enableKeepalive = true
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

func WithDiscovery(d registry.Discovery) ClientOption {
	return func(o *clientOptions) {
		o.discovery = d
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
