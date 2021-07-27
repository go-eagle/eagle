package server

import (
	"github.com/1024casts/snake/internal/routers"
	"github.com/1024casts/snake/internal/service"
	"github.com/1024casts/snake/pkg/conf"
	"github.com/1024casts/snake/pkg/transport/http"
)

// NewHttpServer creates a HTTP server
func NewHttpServer(c *conf.Config, svc *service.Service) *http.Server {
	router := routers.NewRouter()

	var opts []http.ServerOption
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.ReadTimeout != 0 {
		opts = append(opts, http.Timeout(c.Http.ReadTimeout))
	}
	if c.Http.WriteTimeout != 0 {
		opts = append(opts, http.Timeout(c.Http.WriteTimeout))
	}
	srv := http.NewServer(opts...)

	srv.Handler = router
	// NOTE: register svc to http server

	return srv
}
