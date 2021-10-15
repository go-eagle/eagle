package http

import (
	"time"

	"github.com/go-eagle/eagle/pkg/transport"
)

var _ transport.Server = (*Server)(nil)

// ServerOption is HTTP server option
type ServerOption func(*Server)

// WithAddress with server address.
func WithAddress(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// WithReadTimeout with read timeout.
func WithReadTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.readTimeout = timeout
	}
}

// WithWriteTimeout with write timeout.
func WithWriteTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.writeTimeout = timeout
	}
}
