/**
 *    ____          __
 *   / __/__  ___ _/ /_____
 *  _\ \/ _ \/ _ `/  '_/ -_)
 * /___/_//_/\_,_/_/\_\\__/
 *
 * generate by http://patorjk.com/software/taag/#p=display&f=Small%20Slant&t=Snake
 */
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/uber/jaeger-lib/metrics"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"

	"github.com/1024casts/snake/internal/conf"
	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/internal/server/grpc"
	http2 "github.com/1024casts/snake/internal/server/http"
	"github.com/1024casts/snake/internal/service"
	logger "github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/pkg/net/tracing"
	redis2 "github.com/1024casts/snake/pkg/redis"
)

var (
	cfgFile = pflag.StringP("config", "c", "", "snake config file path.")
	Cfg     *conf.Config
	Svc     *service.Service
)

func init() {
	pflag.Parse()
	// init config
	Cfg, err := conf.Init(*cfgFile)
	if err != nil {
		panic(err)
	}
	// init log
	logger.InitLog(&Cfg.Logger)
	// init db
	model.Init(&Cfg.MySQL)
	// init redis
	redis2.Init(&Cfg.Redis)
	// init tracer
	metricsFactory := jprom.New().Namespace(metrics.NSOptions{Name: Cfg.App.Name, Tags: nil})
	_, closer, err := tracing.Init(Cfg.Jaeger.ServiceName, Cfg.Jaeger.Host, metricsFactory)
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	// init service
	Svc := service.New(Cfg)
	_ = Svc
}

// @title snake docs api
// @version 1.0
// @description snake demo

// @host localhost:8080
// @BasePath /v1
func main() {
	gin.SetMode(conf.Conf.App.Mode)
	// init http server
	httpSrv := http2.Init(Svc)
	// init grpc server
	grpcSrv := grpc.Init(Cfg, Svc)
	// init pprof server
	go func() {
		fmt.Printf("Listening and serving PProf HTTP on %s\n", conf.Conf.App.PprofPort)
		if err := http.ListenAndServe(conf.Conf.App.PprofPort, http.DefaultServeMux); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen ListenAndServe for PProf, err: %s", err.Error())
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), conf.Conf.App.CtxDefaultTimeout*time.Second)
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-quit
		log.Printf("Server receive a quit signal: %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Println("Server is exiting")
			// close http server
			if httpSrv != nil {
				if err := httpSrv.Shutdown(ctx); err != nil {
					log.Fatalf("Server shutdown err: %s", err)
				}
			}
			// close grpc server
			if grpcSrv != nil {
				grpcSrv.GracefulStop()
			}
			// close service
			Svc.Close()
			return
		case syscall.SIGHUP:
			// TODO: reload
		default:
			return
		}
	}
}
