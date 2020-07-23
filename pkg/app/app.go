package app

import (
	"github.com/1024casts/snake/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

// App app is singleton
var App *Application

// Application a container for your application.
type Application struct {
	Conf        config.Config
	DB          *gorm.DB
	RedisClient *redis.Client
	Engine      *gin.Engine
	Debug       bool
}

// New create a app
func New() *Application {
	return &Application{
		Conf:        config.Config{},
		DB:          nil,
		RedisClient: nil,
		Engine:      nil,
		Debug:       false,
	}
}

// Run start a app
func (a *Application) Run() {

}
