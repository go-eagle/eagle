package grpc

import (
	"context"
	"net"
	"net/url"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/go-eagle/eagle/pkg/log"
)

// Server is a gRPC server wrapper.
type Server struct {
	*grpc.Server
	ctx      context.Context
	lis      net.Listener
	once     sync.Once
	err      error
	network  string
	address  string
	endpoint *url.URL
	timeout  time.Duration
	inters   []grpc.UnaryServerInterceptor
	grpcOpts []grpc.ServerOption
	health   *health.Server
	log      log.Logger
}

// NewServer creates a gRPC server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
		timeout: 1 * time.Second,
		health:  health.NewServer(),
		log:     log.GetLogger(),
	}
	for _, o := range opts {
		o(srv)
	}
	var ints = []grpc.UnaryServerInterceptor{}
	if len(srv.inters) > 0 {
		ints = append(ints, srv.inters...)
	}
	var grpcOpts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(ints...),
	}
	srv.Server = grpc.NewServer(grpcOpts...)

	// internal register
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	reflection.Register(srv.Server)

	return srv
}

// Start start the gRPC server.
func (s *Server) Start(ctx context.Context) error {
	s.ctx = ctx
	s.log.Infof("[gRPC] server is listening on: %s", s.lis.Addr().String())
	return s.Serve(s.lis)
}

// Stop stop the gRPC server.
func (s *Server) Stop(ctx context.Context) error {
	s.GracefulStop()
	s.log.Info("[gRPC] server is stopping")
	return nil
}
