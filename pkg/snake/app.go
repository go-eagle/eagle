package snake

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"

	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/internal/service"
	"github.com/1024casts/snake/pkg/conf"
	"github.com/1024casts/snake/pkg/log"
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
	Conf        *conf.Config
	DB          *gorm.DB
	RedisClient *redis.Client
	Router      *gin.Engine
	BizService  *service.Service
	RPCServer   *grpc.Server
	Debug       bool
}

// New create a app
func New(cfg *conf.Config) *Application {
	app := new(Application)

	// init log
	conf.InitLog(cfg)

	// init db
	app.DB = model.Init(cfg)

	// init redis
	app.RedisClient = redis2.Init(cfg)

	// init router
	app.Router = gin.Default()

	if cfg.App.RunMode == ModeDebug {
		app.DB.Debug()
		app.Debug = true
	}

	return app
}

// Run start a app
func (a *Application) Run() {
	log.Infof("Start to listening the incoming requests on http address: %s", conf.Conf.App.Addr)
	srv := &http.Server{
		Addr:    conf.Conf.App.Addr,
		Handler: a.Router,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err.Error())
		}
	}()

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
		log.Infof("[snake] Server receive a quit signal: %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("[snake] Server is exiting")

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
