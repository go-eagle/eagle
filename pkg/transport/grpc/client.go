package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc/credentials"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	grpcInsecure "google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
)

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
		unaryClientInterceptor(),
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
	} else {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		cred := credentials.NewTLS(tlsConfig)
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(cred))
	}
	if options.enableGzip {
		dialOpts = append(dialOpts, grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	}
	if options.enableMetric {
		dialOpts = append(dialOpts,
			grpc.WithChainUnaryInterceptor(grpcPrometheus.UnaryClientInterceptor),
			grpc.WithChainStreamInterceptor(grpcPrometheus.StreamClientInterceptor),
		)
	}

	return grpc.DialContext(ctx, options.endpoint, dialOpts...)
}
