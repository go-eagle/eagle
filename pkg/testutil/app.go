// 单元测试工具

package testutil

import (
	"github.com/1024casts/snake/internal/model"
	"github.com/1024casts/snake/pkg/conf"

	"github.com/jinzhu/gorm"
)

// App 结构体，主要是为了方便实例一个app
type App struct {
	DB *gorm.DB
}

// Initialize 初始化
func (app *App) Initialize() {
	// init config
	cfg := "./conf/config.sample.yaml"
	if err := conf.Init(cfg); err != nil {
		panic(err)
	}

	// init db
	model.Init()
	app.DB = model.DB
}
