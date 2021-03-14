// 单元测试工具

package testutil

import (
	"github.com/1024casts/snake/pkg/log"
	"github.com/jinzhu/gorm"

	"github.com/1024casts/snake/internal/conf"
	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/internal/service"
	"github.com/1024casts/snake/pkg/redis"
)

// App 结构体，主要是为了方便实例一个app
type App struct {
	DB *gorm.DB
}

// Initialize 初始化
func (app *App) Initialize() {
	cfg := "../../../../conf/config.sample.yaml"
	if _, err := conf.Init(cfg); err != nil {
		panic(err)
	}

	// init log
	log.InitLog(&conf.Conf.Logger)

	// init db
	model.Init(&conf.Conf.MySQL)

	redis.InitTestRedis()

	// init service
	_ = service.New(conf.Conf)
}
