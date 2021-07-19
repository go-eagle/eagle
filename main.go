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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	_ "go.uber.org/automaxprocs"

	"github.com/1024casts/snake/internal/dao"
	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/internal/server"
	"github.com/1024casts/snake/internal/service"
	"github.com/1024casts/snake/pkg/app"
	"github.com/1024casts/snake/pkg/conf"
	logger "github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/pkg/redis"
	"github.com/1024casts/snake/pkg/trace"
	v "github.com/1024casts/snake/pkg/version"
)

var (
	cfgFile = pflag.StringP("config", "c", "", "snake config file path.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

// @title snake docs api
// @version 1.0
// @description snake demo

// @host localhost:8080
// @BasePath /v1
func main() {
	pflag.Parse()
	if *version {
		ver := v.Get()
		marshaled, err := json.MarshalIndent(&ver, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshaled))
		return
	}

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
	_, err = trace.InitTracerProvider(cfg.Trace.ServiceName, cfg.Trace.Jaeger.CollectorEndpoint)
	if err != nil {
		panic(err)
	}

	// init service
	svc := service.New(cfg, dao.New(model.GetDB()))

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
