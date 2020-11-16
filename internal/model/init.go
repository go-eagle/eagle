package model

import (
	"fmt"

	// MySQL driver.
	"github.com/jinzhu/gorm"
	// GORM MySQL
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/1024casts/snake/pkg/conf"
	"github.com/1024casts/snake/pkg/log"
)

// DB 数据库全局变量
var DB *gorm.DB

// Init 初始化数据库
func Init(cfg *conf.Config) *gorm.DB {
	return openDB(cfg)
}

// openDB 链接数据库，生成数据库实例
func openDB(cfg *conf.Config) *gorm.DB {
	c := cfg.MySQL
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		c.UserName,
		c.Password,
		c.Addr,
		c.Name,
		true,
		//"Asia/Shanghai"),
		"Local")

	db, err := gorm.Open("mysql", config)
	if err != nil {
		log.Panicf("Database connection failed. Database name: %s, err: %+v", c.Name, err)
	}

	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	// set for db connection
	db.LogMode(c.ShowLog)
	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxOpenConns(c.MaxOpenConn)
	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	db.DB().SetMaxIdleConns(c.MaxIdleConn)
	db.DB().SetConnMaxLifetime(c.ConnMaxLifeTime)

	DB = db

	return db
}

// GetDB 返回默认的数据库
func GetDB() *gorm.DB {
	return DB
}
