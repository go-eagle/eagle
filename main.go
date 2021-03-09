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
	"github.com/1024casts/snake/config"
	snake "github.com/1024casts/snake/pkg"
	"github.com/1024casts/snake/pkg/net/tracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/uber/jaeger-lib/metrics"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"
	"google.golang.org/grpc"

	"github.com/1024casts/snake/app/api"
	rpc "github.com/1024casts/snake/internal/server"
	"github.com/1024casts/snake/internal/service"
	routers "github.com/1024casts/snake/router"
)

var (
	cfgFile = pflag.StringP("config", "c", "", "snake config file path.")
)

// @title snake docs api
// @version 1.0
// @description snake demo

// @contact.name 1024casts/snake
// @contact.url http://www.swagger.io/support
// @contact.email

// @host localhost:8080
// @BasePath /v1
func main() {
	pflag.Parse()

	// init config
	cfg, err := config.Init(*cfgFile)
	if err != nil {
		panic(err)
	}

	// Set gin mode.
	gin.SetMode(cfg.App.Mode)

	// init app
	app := snake.New(cfg)
	snake.App = app

	// init db tracing plugin
	//app.DB.Use(gormopentracing.New())

	// Create the Gin engine.
	router := app.Router

	// HealthCheck 健康检查路由
	router.GET("/health", api.HealthCheck)
	// metrics router 可以在 prometheus 中进行监控
	// 通过 grafana 可视化查看 prometheus 的监控数据，使用插件6671查看
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API Routes.
	routers.Load(router)
	// WEB Routes
	routers.LoadWebRouter(router)

	// init tracer
	metricsFactory := jprom.New().Namespace(metrics.NSOptions{Name: cfg.App.Name, Tags: nil})
	tracer, closer := tracing.Init(cfg.App.Name, metricsFactory)
	defer closer.Close()

	// set into opentracing
	opentracing.SetGlobalTracer(tracer)

	// init service
	svc := service.New(cfg, tracer)

	// set global service
	service.Svc = svc
	snake.App.BizService = svc

	// start grpc server
	var rpcSrv *grpc.Server
	go func() {
		rpcSrv = rpc.New(cfg, svc)
		snake.App.RPCServer = rpcSrv
	}()

	// here register to service discovery

	// start server
	app.Run()
}
