package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	MetadataClientAPIVersionKey = "client-api-version"
)

// unaryServerInterceptor server unary interceptor
func unaryServerInterceptor(s *Server) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// preprocess stage
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			ver, vs := "unknown", md.Get(MetadataClientAPIVersionKey)
			if len(vs) > 0 {
				ver = vs[0]
			}
			clientRequests.WithLabelValues("unary", ver).Inc()
		}

		return handler(ctx, req)
	}
}

// unaryClientInterceptor client unary interceptor
func unaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// preprocess stage

		// call remote method
		err := invoker(ctx, method, req, reply, cc, opts...)

		// postprocess stage

		return err
	}
}
