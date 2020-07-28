package snake

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	redis2 "github.com/1024casts/snake/pkg/redis"

	"github.com/1024casts/snake/pkg/schedule"

	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/pkg/config"
	"github.com/1024casts/snake/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

const (
	// ModeDebug debug mode
	ModeDebug string = "debug"
	// ModeRelease release mode
	ModeRelease string = "release"
	// ModeTest test mode
	ModeTest string = "test"
)

// App is singleton
var App *Application

// Application a container for your application.
type Application struct {
	Conf        *config.Config
	DB          *gorm.DB
	RedisClient *redis.Client
	Router      *gin.Engine
	Debug       bool
}

// New create a app
func New(conf *config.Config) *Application {
	app := new(Application)

	// init db
	app.DB = model.Init()

	// init redis
	app.RedisClient = redis2.Init()

	// init router
	app.Router = gin.Default()

	// init log
	config.InitLog()

	// init schedule
	schedule.Init()

	if viper.GetString("run_mode") == ModeDebug {
		app.Debug = true
	}
	App = app

	return app
}

// Run start a app
func (a *Application) Run() {
	log.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("app.addr"))
	srv := &http.Server{
		Addr:    viper.GetString("app.addr"),
		Handler: a.Router,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err.Error())
		}
	}()

	gracefulStop(srv)
}

// gracefulStop 优雅退出
// 等待中断信号以超时 5 秒正常关闭服务器
// 官方说明：https://github.com/gin-gonic/gin#graceful-restart-or-stop
func gracefulStop(srv *http.Server) {
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
	default:
	}
	log.Info("Server exiting")
}
