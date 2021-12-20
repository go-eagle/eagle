package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	MetadataClientAPIVersionKey = "client-api-version"
)

func newUnaryInterceptor(s *Server) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
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
