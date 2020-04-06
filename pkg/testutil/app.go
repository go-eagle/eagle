package testutil

import (
	"github.com/1024casts/snake/config"
	"github.com/1024casts/snake/model"
	"github.com/jinzhu/gorm"
)

type App struct {
	DB *gorm.DB
}

func (app *App) Initialize() {
	// init config
	if err := config.Init("../../conf/config.sample.yaml"); err != nil {
		panic(err)
	}

	// init db
	model.DB.Init()
	app.DB = model.DB.Self
}
