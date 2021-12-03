package grpc

import (
	"context"
	"log"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/go-eagle/eagle/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Server is a gRPC server wrapper.
// nolint
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
	}
	for _, o := range opts {
		o(srv)
	}
	var inters []grpc.UnaryServerInterceptor
	if len(srv.inters) > 0 {
		inters = append(inters, srv.inters...)
	}
	var grpcOpts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(inters...),
	}
	srv.Server = grpc.NewServer(grpcOpts...)

	// internal register
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	reflection.Register(srv.Server)

	return srv
}

// Endpoint return a real address to registry endpoint.
// examples:
//   grpc://127.0.0.1:9090
func (s *Server) Endpoint() (*url.URL, error) {
	addr, err := utils.Extract(s.address, s.lis)
	if err != nil {
		return nil, err
	}
	s.endpoint = &url.URL{Scheme: "grpc", Host: addr}
	return s.endpoint, nil
}

// Start start the gRPC server.
func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	s.lis = lis

	if _, err := s.Endpoint(); err != nil {
		return err
	}

	s.ctx = ctx
	log.Printf("[gRPC] server is listening on: %s", s.lis.Addr().String())
	return s.Serve(s.lis)
}

// Stop stop the gRPC server.
func (s *Server) Stop(ctx context.Context) error {
	s.GracefulStop()
	log.Printf("[gRPC] server is stopping")
	return nil
}
