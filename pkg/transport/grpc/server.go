package grpc

import (
	"context"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/1024casts/snake/pkg/log"

	"google.golang.org/grpc"
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
	log      log.Logger
}

// NewServer creates a gRPC server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
		timeout: 1 * time.Second,
		log:     log.GetLogger(),
	}
	for _, o := range opts {
		o(srv)
	}
	srv.Server = grpc.NewServer()
	return srv
}

// Start start the gRPC server.
func (s *Server) Start(ctx context.Context) error {
	s.log.Infof("[gRPC] server is listening on: %s", s.lis.Addr().String())
	return s.Serve(s.lis)
}

// Stop stop the gRPC server.
func (s *Server) Stop(ctx context.Context) error {
	s.GracefulStop()
	s.log.Info("[gRPC] server is stopping")
	return nil
}
