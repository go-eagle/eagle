package snake

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/1024casts/snake/config"
	"github.com/1024casts/snake/internal/service"
	"github.com/1024casts/snake/pkg/database/orm"

	//"github.com/1024casts/snake/pkg/log"
	redis2 "github.com/1024casts/snake/pkg/redis"
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
	BizService  *service.Service
	RPCServer   *grpc.Server
	Debug       bool
}

// New create a app
func New(cfg *config.Config) *Application {
	app := new(Application)

	// init log
	config.InitLog(cfg)

	// init db
	app.DB = orm.NewMySQL(&cfg.MySQL)

	// init redis
	app.RedisClient = redis2.Init(cfg)

	// init router
	app.Router = gin.Default()

	if cfg.App.Mode == ModeDebug {
		app.DB.Debug()
		app.Debug = true
	}

	return app
}

// Run start a app
func (a *Application) Run() {
	fmt.Printf("Listening and serving HTTP on %s\n", config.Conf.App.Port)
	srv := &http.Server{
		Addr:    config.Conf.App.Port,
		Handler: a.Router,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen, err: %s", err.Error())
		}
	}()

	// block process
	a.GracefulStop(srv)
}

// GracefulStop 优雅退出
func (a *Application) GracefulStop(httpSrv *http.Server) {
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-quit
		log.Printf("[snake] Server receive a quit signal: %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Println("[snake] Server is exiting")

			// close rpc
			if a.RPCServer != nil {
				a.RPCServer.GracefulStop()
			}

			// close biz svc
			a.BizService.Close()

			// close http server
			if httpSrv != nil {
				if err := httpSrv.Shutdown(context.Background()); err != nil {
					log.Fatalf("[snake] Server shutdown err: %s", err)
				}
			}
			return
		case syscall.SIGHUP:
			//todo: reload
		default:
			return
		}
	}
}
