package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/1024casts/snake/internal/conf"
	"github.com/1024casts/snake/internal/routers"
	"github.com/1024casts/snake/internal/service"
)

// NewHttpServer creates a HTTP server
func NewHttpServer(s *service.Service) *http.Server {
	router := routers.NewRouter()
	srv := &http.Server{
		Addr:         conf.Conf.Http.Addr,
		Handler:      router,
		ReadTimeout:  time.Second * conf.Conf.Http.ReadTimeout,
		WriteTimeout: time.Second * conf.Conf.Http.WriteTimeout,
	}

	fmt.Printf("Listening and serving HTTP on %s\n", conf.Conf.Http.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe, err: %s", err.Error())
	}

	return srv
}
