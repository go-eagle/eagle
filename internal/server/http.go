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
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	srv := http.NewServer(opts...)
	srv.Handler = router

	return srv
}
