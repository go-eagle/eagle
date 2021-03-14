package http

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/1024casts/snake/internal/conf"
	"github.com/1024casts/snake/internal/routers"
	"github.com/1024casts/snake/internal/service"
)

var (
	UserSvc *service.Service
)

func Init(s *service.Service) *http.Server {
	UserSvc = s
	router := routers.NewRouter()
	srv := &http.Server{
		Addr:         conf.Conf.App.Port,
		Handler:      router,
		ReadTimeout:  time.Second * conf.Conf.App.ReadTimeout,
		WriteTimeout: time.Second * conf.Conf.App.WriteTimeout,
	}

	fmt.Printf("Listening and serving HTTP on %s\n", conf.Conf.App.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe, err: %s", err.Error())
	}

	return srv
}
