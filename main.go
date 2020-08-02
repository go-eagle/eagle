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
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/1024casts/snake/handler"
	"github.com/1024casts/snake/pkg/conf"
	"github.com/1024casts/snake/pkg/snake"
	v "github.com/1024casts/snake/pkg/version"
	routers "github.com/1024casts/snake/router"
)

var (
	version = pflag.BoolP("version", "v", false, "show version info.")
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
	if err := conf.Init(); err != nil {
		panic(err)
	}

	// init app
	snake.App = snake.New()

	// Set gin mode.
	gin.SetMode(snake.ModeRelease)
	if viper.GetString("app.run_mode") == snake.ModeDebug {
		gin.SetMode(snake.ModeDebug)
		snake.App.DB.Debug()
	}

	// Create the Gin engine.
	router := snake.App.Router

	// HealthCheck 健康检查路由
	router.GET("/health", handler.HealthCheck)
	// metrics router 可以在 prometheus 中进行监控
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API Routes.
	routers.Load(router)

	// start server
	snake.App.Run()
}
