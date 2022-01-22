package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	grpcInsecure "google.golang.org/grpc/credentials/insecure"
)

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

// clientOptions define gRPc client options.
type clientOptions struct {
	endpoint     string
	timeout      time.Duration
	inters       []grpc.UnaryClientInterceptor
	dialOpts     []grpc.DialOption
	balancerName string
}

// Dial
func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	// default client options
	options := clientOptions{
		timeout:      2000 * time.Millisecond,
		balancerName: roundrobin.Name,
	}
	for _, opt := range opts {
		opt(&options)
	}

	// merge inters
	inters := []grpc.UnaryClientInterceptor{
		// here add default unary client interceptor
	}
	if len(options.inters) > 0 {
		inters = append(inters, options.inters...)
	}

	// default dial option
	dialOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy": "%s"}`, options.balancerName)),
		grpc.WithChainUnaryInterceptor(inters...),
	}
	if len(options.dialOpts) > 0 {
		dialOpts = append(dialOpts, options.dialOpts...)
	}
	if insecure {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(grpcInsecure.NewCredentials()))
	}

	return grpc.DialContext(ctx, options.endpoint, dialOpts...)
}
