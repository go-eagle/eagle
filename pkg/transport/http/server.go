package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/1024casts/snake/pkg/log"
)

type Server struct {
	*http.Server
	lis     net.Listener
	network string
	address string
	timeout time.Duration
	log     log.Logger
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":8080",
		timeout: time.Second,
		log:     log.GetLogger(),
	}
	for _, o := range opts {
		o(srv)
	}
	// NOTE: must set server
	srv.Server = &http.Server{Handler: srv}
	return srv
}

// ServeHTTP should write reply headers and data to the ResponseWriter and then return.
func (s *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

}

func (s *Server) Start() error {
	lis, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	s.lis = lis
	s.log.Infof("http server is listening on: %s", lis.Addr().String())
	if err := s.Serve(lis); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	s.log.Info("http server is stopping")
	return s.Shutdown(context.Background())
}
