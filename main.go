/**
 *
 *    ____          __
 *   / __/__ ____ _/ /__
 *  / _// _ `/ _ `/ / -_)
 * /___/\_,_/\_, /_/\__/
 *         /___/
 *
 *
 * generate by http://patorjk.com/software/taag/#p=display&f=Small%20Slant&t=Eagle
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

	"github.com/go-eagle/eagle/internal/dao"
	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/internal/server"
	"github.com/go-eagle/eagle/internal/service"
	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/conf"
	logger "github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/redis"
	"github.com/go-eagle/eagle/pkg/trace"
	v "github.com/go-eagle/eagle/pkg/version"
)

var (
	cfgFile = pflag.StringP("config", "c", "", "eagle config file path.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

// @title eagle docs api
// @version 1.0
// @description eagle demo

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
	model.Init(&cfg.ORM)
	// init redis
	redis.Init(&cfg.Redis)
	// init tracer
	_, err = trace.InitTracerProvider(cfg.Trace.ServiceName, cfg.Trace.Jaeger.CollectorEndpoint)
	if err != nil {
		panic(err)
	}

	// init service
	svc := service.New(cfg, dao.New(cfg, model.GetDB()))

	gin.SetMode(conf.Conf.App.Mode)

	// init pprof server
	go func() {
		fmt.Printf("Listening and serving PProf HTTP on %s\n", conf.Conf.App.PprofPort)
		if err := http.ListenAndServe(conf.Conf.App.PprofPort, http.DefaultServeMux); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen ListenAndServe for PProf, err: %s", err.Error())
		}
	}()

	app := eagle.New(cfg,
		eagle.WithName(cfg.App.Name),
		eagle.WithVersion(cfg.App.Version),
		eagle.WithLogger(logger.GetLogger()),
		eagle.Server(
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
