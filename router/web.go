package routers

import (
	"html/template"
	"time"

	gintemplate "github.com/foolin/gin-template"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/1024casts/snake/app/web"
	webUser "github.com/1024casts/snake/app/web/user"
	"github.com/1024casts/snake/pkg/flash"
	"github.com/1024casts/snake/pkg/log"
)

// LoadWebRouter loads the middlewares, routes, handlers.
func LoadWebRouter(g *gin.Engine) *gin.Engine {
	router := g

	// Middlewares.

	// 404 Handler.
	router.NoRoute(func(c *gin.Context) {
		web.Error404(c)
	})
	router.NoMethod(func(c *gin.Context) {
		web.Error404(c)
	})

	router.Use(static.Serve("/static", static.LocalFile(viper.GetString("static"), false)))
	router.Use(static.Serve("/uploads/avatar", static.LocalFile(viper.GetString("avatar"), false)))
	router.Use(static.Serve("/uploads/images", static.LocalFile(viper.GetString("images"), false)))

	//new template engine
	router.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		Root:      "internal/templates",
		Extension: ".html",
		Master:    "layouts/master",
		Partials:  []string{},
		Funcs: template.FuncMap{
			// 判断是否是当前链接
			"isActive": func(ctx *gin.Context, currentUri string) string {
				if ctx.Request.RequestURI == currentUri {
					return "is-active"
				} else {
					return ""
				}
			},
			// 全局消息
			"flashMessage": func(ctx *gin.Context) string {
				errorMessage, err := flash.GetMessage(ctx.Writer, ctx.Request)
				if err != nil {
					log.Warnf("[router] get flash message err: %v", err)
					return ""
				}
				return string(errorMessage)
			},
			"hasFlash": func(ctx *gin.Context) bool {
				return flash.HasFlash(ctx.Request)
			},
			"copy": func() string {
				return time.Now().Format("2006")
			},
		},
		DisableCache: true,
	})

	router.GET("/", web.Index)

	// login
	router.GET("/login", webUser.GetLogin)
	router.POST("/login", webUser.DoLogin)
	router.GET("/logout", webUser.Logout)

	// register
	router.GET("/register", webUser.GetRegister)
	router.POST("/register", webUser.DoRegister)

	return router
}
