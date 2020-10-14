package model

import (
	"fmt"
	"time"

	"github.com/spf13/viper"

	// MySQL driver.
	"github.com/jinzhu/gorm"
	// GORM MySQL
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/1024casts/snake/pkg/log"
)

// DB 数据库全局变量
var DB *gorm.DB

// Init 初始化数据库
func Init() *gorm.DB {
	return openDB(viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.addr"),
		viper.GetString("mysql.name"))
}

// openDB 链接数据库，生成数据库实例
func openDB(username, password, addr, name string) *gorm.DB {
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		username,
		password,
		addr,
		name,
		true,
		//"Asia/Shanghai"),
		"Local")

	db, err := gorm.Open("mysql", config)
	if err != nil {
		log.Errorf("Database connection failed. Database name: %s, err: %+v", name, err)
		panic(err)
	}

	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	// set for db connection
	db.LogMode(viper.GetBool("mysql.show_log"))
	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxOpenConns(viper.GetInt("mysql.max_open_conn"))
	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	db.DB().SetMaxIdleConns(viper.GetInt("mysql.max_idle_conn"))
	db.DB().SetConnMaxLifetime(time.Minute * viper.GetDuration("mysql.conn_max_life_time"))

	DB = db

	return db
}

// GetDB 返回默认的数据库
func GetDB() *gorm.DB {
	return DB
}
