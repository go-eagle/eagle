package routers

import (
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger" //nolint: goimports
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/go-eagle/eagle/internal/handler/v1/user"
	mw "github.com/go-eagle/eagle/internal/middleware"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/middleware"
)

// NewRouter loads the middlewares, routes, handlers.
func NewRouter() *gin.Engine {
	g := gin.New()
	// 使用中间件
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(middleware.Logging())
	g.Use(middleware.RequestID())
	g.Use(middleware.Metrics(app.Conf.Name))
	g.Use(middleware.Tracing(app.Conf.Name))
	g.Use(middleware.Timeout(3 * time.Second))
	g.Use(mw.Translations())

	// load web router
	LoadWebRouter(g)

	// 404 Handler.
	g.NoRoute(app.RouteNotFound)
	g.NoMethod(app.RouteNotFound)

	// swagger api docs
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// pprof router 性能分析路由
	// 默认关闭，开发环境下可以打开
	// 访问方式: HOST/debug/pprof
	// 通过 HOST/debug/pprof/profile 生成profile
	// 查看分析图 go tool pprof -http=:5000 profile
	// see: https://github.com/gin-contrib/pprof
	if app.Conf.EnablePprof {
		pprof.Register(g)
	}

	// HealthCheck 健康检查路由
	g.GET("/health", app.HealthCheck)
	// metrics router 可以在 prometheus 中进行监控
	// 通过 grafana 可视化查看 prometheus 的监控数据，使用插件6671查看
	g.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// v1 router
	apiV1 := g.Group("/v1")
	apiV1.Use()
	{
		// 认证相关路由
		apiV1.POST("/register", user.Register)
		apiV1.POST("/login", user.Login)
		apiV1.POST("/login/phone", user.PhoneLogin)
		apiV1.GET("/vcode", user.VCode)

		// 用户
		apiV1.GET("/users/:id", user.Get)
		apiV1.Use(middleware.Auth())
		{
			apiV1.PUT("/users/:id", user.Update)
			apiV1.POST("/users/follow", user.Follow)
			apiV1.POST("/users/unfollow", user.Unfollow)
			apiV1.GET("/users/:id/following", user.FollowList)
			apiV1.GET("/users/:id/followers", user.FollowerList)
		}
	}

	return g
}
