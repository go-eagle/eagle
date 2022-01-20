package grpc

import (
	"context"
	"log"
	"net"
	"net/url"
	"sync"
	"time"

	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/health"
	healthPb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/go-eagle/eagle/pkg/utils"
)

// ServerOption is gRPC server option.
type ServerOption func(o *Server)

// Network with server network.
func Network(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

// Address with server address.
func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// Timeout with server timeout.
func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// UnaryInterceptor returns a ServerOption that sets the UnaryServerInterceptor for the server.
func UnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.inters = in
	}
}

// Options with grpc options.
func Options(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}

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

	// EnableTracer enables distributed tracing using OpenTelemetry protocol
	EnableTracing bool
	// TracerOptions are options for OpenTelemetry gRPC interceptor.
	TracerOptions []otelgrpc.Option
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
	// Unary
	chainUnaryInterceptors := []grpc.UnaryServerInterceptor{
		newUnaryInterceptor(srv),
		grpcPrometheus.UnaryServerInterceptor,
		grpcRecovery.UnaryServerInterceptor(),
	}
	if len(srv.inters) > 0 {
		chainUnaryInterceptors = append(chainUnaryInterceptors, srv.inters...)
	}

	// stream
	chainStreamInterceptors := []grpc.StreamServerInterceptor{
		grpcPrometheus.StreamServerInterceptor,
		grpcRecovery.StreamServerInterceptor(),
	}

	// enable tracing
	if srv.EnableTracing {
		chainUnaryInterceptors = append(chainUnaryInterceptors, otelgrpc.UnaryServerInterceptor(srv.TracerOptions...))
		chainStreamInterceptors = append(chainStreamInterceptors, otelgrpc.StreamServerInterceptor(srv.TracerOptions...))
	}

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(chainUnaryInterceptors...),
		grpc.ChainStreamInterceptor(chainStreamInterceptors...),
	}
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}

	grpcServer := grpc.NewServer(grpcOpts...)

	// see https://github.com/grpc/grpc/blob/master/doc/health-checking.md for more
	srv.health.SetServingStatus("", healthPb.HealthCheckResponse_SERVING)
	healthPb.RegisterHealthServer(grpcServer, srv.health)
	reflection.Register(grpcServer)

	// set zero values for metrics registered for this grpc server
	grpcPrometheus.Register(grpcServer)

	srv.Server = grpcServer

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
