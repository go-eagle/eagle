package server

import (
	"github.com/go-eagle/eagle/internal/routers"
	"github.com/go-eagle/eagle/pkg/conf"
	"github.com/go-eagle/eagle/pkg/transport/http"
)

// NewHTTPServer creates a HTTP server
func NewHTTPServer(c *conf.Config) *http.Server {
	router := routers.NewRouter()

	srv := http.NewServer(
		http.WithAddress(c.HTTP.Addr),
		http.WithReadTimeout(c.HTTP.ReadTimeout),
		http.WithWriteTimeout(c.HTTP.WriteTimeout),
	)

	srv.Handler = router
	// NOTE: register svc to http server

	return srv
}
