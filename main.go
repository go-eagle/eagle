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
	"fmt"
	"log"
	"net/http"

	"github.com/1024casts/snake/pkg/app"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/uber/jaeger-lib/metrics"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"

	"github.com/1024casts/snake/internal/conf"
	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/internal/server"
	"github.com/1024casts/snake/internal/service"
	logger "github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/pkg/net/tracing"
	"github.com/1024casts/snake/pkg/redis"
)

var (
	cfgFile = pflag.StringP("config", "c", "", "snake config file path.")
)

// @title snake docs api
// @version 1.0
// @description snake demo

// @host localhost:8080
// @BasePath /v1
func main() {
	pflag.Parse()
	// init config
	cfg, err := conf.Init(*cfgFile)
	if err != nil {
		panic(err)
	}
	logger.Init(&cfg.Logger)
	// init db
	model.Init(&cfg.MySQL)
	// init redis
	redis.Init(&cfg.Redis)
	// init tracer
	metricsFactory := jprom.New().Namespace(metrics.NSOptions{Name: cfg.App.Name, Tags: nil})
	_, closer, err := tracing.Init(cfg.Jaeger.ServiceName, cfg.Jaeger.Host, metricsFactory)
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	// init service
	svc := service.New(cfg)

	gin.SetMode(conf.Conf.App.Mode)

	// init pprof server
	go func() {
		fmt.Printf("Listening and serving PProf HTTP on %s\n", conf.Conf.App.PprofPort)
		if err := http.ListenAndServe(conf.Conf.App.PprofPort, http.DefaultServeMux); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen ListenAndServe for PProf, err: %s", err.Error())
		}
	}()

	app := app.New(cfg,
		app.WithName(cfg.App.Name),
		app.WithVersion(cfg.App.Version),
		app.WithLogger(logger.GetLogger()),
		app.Server(
			// init http server
			server.NewHttpServer(conf.Conf, svc),
			// init grpc server
			//grpcSrv := server.NewGRPCServer(svc)
		),
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
