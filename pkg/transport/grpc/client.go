package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials"
	grpcInsecure "google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/keepalive"

	logger "github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/transport/grpc/resolver/discovery"
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
		timeout:         2000 * time.Millisecond,
		balancerName:    roundrobin.Name,
		enableGzip:      true,
		enableMetric:    true,
		disableRetry:    false,
		NumRetries:      2,
		enableKeepalive: true,
		kp: keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             time.Second,
			PermitWithoutStream: false,
		},
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
	// service discovery
	if options.discovery != nil {
		dialOpts = append(dialOpts, grpc.WithResolvers(discovery.NewBuilder(
			options.discovery, discovery.WithInsecure(insecure))))
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
	if options.enableKeepalive {
		kp := keepalive.ClientParameters{
			Time:                options.kp.Time,
			Timeout:             options.kp.Timeout,
			PermitWithoutStream: options.kp.PermitWithoutStream,
		}
		dialOpts = append(dialOpts, grpc.WithKeepaliveParams(kp))
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
	// enable tracing
	if options.enableTracing {
		dialOpts = append(dialOpts,
			grpc.WithChainUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
			grpc.WithChainStreamInterceptor(otelgrpc.StreamClientInterceptor()),
		)
	}
	if options.enableLog {
		dialOpts = append(dialOpts,
			grpc.WithChainUnaryInterceptor(grpcZap.UnaryClientInterceptor(logger.GetZapLogger())),
			grpc.WithChainStreamInterceptor(grpcZap.StreamClientInterceptor(logger.GetZapLogger())),
		)
	}
	if !options.disableRetry {
		dialOpts = append(dialOpts,
			grpc.WithDefaultServiceConfig(getRetryPolicy(options.balancerName, options.NumRetries)),
		)
	}

	return grpc.DialContext(ctx, options.endpoint, dialOpts...)
}

func getRetryPolicy(balancerName string, numRetries int) string {
	retryPolicy := `{
		"loadBalancingPolicy": "%s",
		"methodConfig": [{
		  "retryPolicy": {
			  "MaxAttempts": %d,
			  "InitialBackoff": ".01s",
			  "MaxBackoff": ".01s",
			  "BackoffMultiplier": 1.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`
	return fmt.Sprintf(retryPolicy, balancerName, numRetries)
}
