package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/go-eagle/eagle/pkg/log"
)

// Server http server struct
type Server struct {
	*http.Server
	lis          net.Listener
	network      string
	address      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	log          log.Logger
}

// defaultServer return a default config server
func defaultServer() *Server {
	return &Server{
		network:      "tcp",
		address:      ":8080",
		readTimeout:  time.Second,
		writeTimeout: time.Second,
		log:          log.GetLogger(),
	}
}

// NewServer create a server
func NewServer(opts ...ServerOption) *Server {
	srv := defaultServer()
	// apply options
	for _, o := range opts {
		o(srv)
	}
	// NOTE: must set server
	srv.Server = &http.Server{
		ReadTimeout:  srv.readTimeout,
		WriteTimeout: srv.writeTimeout,
		Handler:      srv,
	}
	return srv
}

// ServeHTTP should write reply headers and data to the ResponseWriter and then return.
func (s *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	s.ServeHTTP(resp, req)
}

// Start start a server
func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	s.lis = lis
	s.log.Infof("[HTTP] server is listening on: %s", lis.Addr().String())
	if err := s.Serve(lis); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop stop server
func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("[HTTP] server is stopping")
	return s.Shutdown(ctx)
}
