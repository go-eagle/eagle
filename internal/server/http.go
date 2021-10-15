package server

import (
	"github.com/go-eagle/eagle/internal/routers"
	"github.com/go-eagle/eagle/pkg/conf"
	"github.com/go-eagle/eagle/pkg/transport/http"
)

// NewHttpServer creates a HTTP server
func NewHttpServer(c *conf.Config) *http.Server {
	router := routers.NewRouter()

	srv := http.NewServer(
		http.WithAddress(c.Http.Addr),
		http.WithReadTimeout(c.Http.ReadTimeout),
		http.WithWriteTimeout(c.Http.WriteTimeout),
	)

	srv.Handler = router
	// NOTE: register svc to http server

	return srv
}
