package model

import (
	"github.com/1024casts/snake/config"
	"gorm.io/gorm"

	"github.com/1024casts/snake/pkg/database/orm"
)

// DB 数据库全局变量
var DB *gorm.DB

// Init 初始化数据库
func Init(cfg *config.Config) *gorm.DB {
	DB = orm.NewMySQL(&cfg.MySQL)
	return DB
}

// GetDB 返回默认的数据库
func GetDB() *gorm.DB {
	return DB
}
