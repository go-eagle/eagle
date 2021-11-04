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

	"github.com/go-eagle/eagle/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	etcdclient "go.etcd.io/etcd/client/v3"
	_ "go.uber.org/automaxprocs"

	"github.com/go-eagle/eagle/internal/dao"
	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/internal/server"
	"github.com/go-eagle/eagle/internal/service"
	eagle "github.com/go-eagle/eagle/pkg/app"
	logger "github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/redis"
	"github.com/go-eagle/eagle/pkg/registry/etcd"
	"github.com/go-eagle/eagle/pkg/trace"
	v "github.com/go-eagle/eagle/pkg/version"
)

var (
	cfgDir  = pflag.StringP("config dir", "d", "config", "eagle config file dir.")
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
	c := config.New(config.WithConfigDir(*cfgDir))
	var cfg config.AppConfig
	if err := c.Load("app", &cfg); err != nil {
		panic(err)
	}
	// set global
	config.App = cfg

	// -------------- init resource -------------
	logger.Init()
	// init db
	model.Init()
	// init redis
	redis.Init()
	// init tracer
	if cfg.EnableTrace {
		var traceCfg trace.Config
		err := config.Conf.Load("trace", &traceCfg)
		_, err = trace.InitTracerProvider(traceCfg.ServiceName, traceCfg.CollectorEndpoint)
		if err != nil {
			panic(err)
		}
	}

	// init service
	service.Svc = service.New(dao.New(model.GetDB()))

	gin.SetMode(cfg.Mode)

	// init pprof server
	go func() {
		fmt.Printf("Listening and serving PProf HTTP on %s\n", cfg.PprofPort)
		if err := http.ListenAndServe(cfg.PprofPort, http.DefaultServeMux); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen ListenAndServe for PProf, err: %s", err.Error())
		}
	}()

	// create a etcd register
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		log.Fatal(err)
	}
	r := etcd.New(client)

	// start app
	app := eagle.New(
		eagle.WithName(cfg.Name),
		eagle.WithVersion(cfg.Version),
		eagle.WithLogger(logger.GetLogger()),
		eagle.WithServer(
			// init http server
			server.NewHTTPServer(&cfg.HTTP),
		),
		eagle.WithRegistry(r),
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
