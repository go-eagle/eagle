package snake

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/1024casts/snake/pkg/conf"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/1024casts/snake/internal/service"
	logger "github.com/1024casts/snake/pkg/log"
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
	logger.InitLog(cfg)

	// init router
	app.Router = gin.Default()

	return app
}

// Run start a app
func (a *Application) Run() {
	srv := &http.Server{
		Addr:         conf.Conf.App.Port,
		ReadTimeout:  time.Second * conf.Conf.App.ReadTimeout,
		WriteTimeout: time.Second * conf.Conf.App.WriteTimeout,
		Handler:      a.Router,
	}

	go func() {
		fmt.Printf("Listening and serving HTTP on %s\n", conf.Conf.App.Port)
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen, err: %s", err.Error())
		}
	}()

	go func() {
		fmt.Printf("Listening and serving Pprof HTTP on %s\n", conf.Conf.App.PprofPort)
		// service connections
		if err := http.ListenAndServe(conf.Conf.App.PprofPort, http.DefaultServeMux); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen Pprof, err: %s", err.Error())
		}
	}()

	ctx, shutdown := context.WithTimeout(context.Background(), conf.Conf.App.CtxDefaultTimeout*time.Second)
	defer shutdown()

	// block process
	a.GracefulStop(ctx, srv)
}

// GracefulStop 优雅退出
func (a *Application) GracefulStop(ctx context.Context, httpSrv *http.Server) {
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
				if err := httpSrv.Shutdown(ctx); err != nil {
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
