package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/1024casts/snake/model"
	"github.com/1024casts/snake/pkg/redis"

	"github.com/1024casts/snake/config"
	v "github.com/1024casts/snake/pkg/version"
	"github.com/1024casts/snake/router"
	"github.com/1024casts/snake/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg     = pflag.StringP("config", "c", "", "snake config file path.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

// @title snake Example API
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
		v := v.Get()
		marshaled, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshaled))
		return
	}

	// init config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// init db
	model.DB.Init()
	defer model.DB.Close()

	// init redis
	redis.Init()

	// Set gin mode.
	gin.SetMode(viper.GetString("runmode"))

	// Create the Gin engine.
	g := gin.New()

	// Routes.
	router.Load(
		// Cores.
		g,

		// Middlwares.
		middleware.Logging(),
		middleware.RequestId(),
	)

	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()

	// Start to listening the incoming requests.
	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")
	if cert != "" && key != "" {
		go func() {
			log.Infof("Start to listening the incoming requests on https address: %s", viper.GetString("tls.addr"))
			log.Info(http.ListenAndServeTLS(viper.GetString("tls.addr"), cert, key, g).Error())
		}()
	}

	log.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	srv := &http.Server{
		Addr:    viper.GetString("addr"),
		Handler: g,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(err, "listen: %s", err.Error())
		}
		log.Info(srv.ListenAndServe().Error())
	}()

	// 等待中断信号以超时 5 秒正常关闭服务器
	// 官方说明：https://github.com/gin-gonic/gin#graceful-restart-or-stop
	quit := make(chan os.Signal)
	// kill 命令发送信号 syscall.SIGTERM
	// kill -2 命令发送信号 syscall.SIGINT
	// kill -9 命令发送信号 syscall.SIGKILL
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// 5 秒后捕获 ctx.Done() 信号
	select {
	case <-ctx.Done():
		log.Info("timeout of 5 seconds.")
	}
	log.Info("Server exiting")
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("cannot connect to the router")
}
