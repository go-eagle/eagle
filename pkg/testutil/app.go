// 单元测试工具

package testutil

import (
	"github.com/1024casts/snake/config"
	"github.com/1024casts/snake/model"

	"github.com/jinzhu/gorm"
)

// App 结构体，主要是为了方便实例一个app
type App struct {
	DB *gorm.DB
}

// Initialize 初始化
func (app *App) Initialize() {
	// init config
	if err := config.Init("../../conf/config.sample.yaml"); err != nil {
		panic(err)
	}

	// init db
	model.DB.Init()
	app.DB = model.DB.Default
}
