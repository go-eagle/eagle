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

// Database 定义现有的数据库
type Database struct {
	Default *gorm.DB
	Docker  *gorm.DB
}

// DB 数据库全局变量
var DB *Database

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
	}

	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	// set for db connection
	setupDB(db)

	return db
}

// setupDB 配置数据库
func setupDB(db *gorm.DB) {
	db.LogMode(viper.GetBool("gorm.show_log"))
	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxOpenConns(viper.GetInt("grom.max_open_conn"))
	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	db.DB().SetMaxIdleConns(viper.GetInt("grom.max_idle_conn"))
	db.DB().SetConnMaxLifetime(time.Minute * viper.GetDuration("grom.conn_max_lift_time"))
}

// GetDB 返回默认的数据库
func GetDB() *gorm.DB {
	return DB.Default
}

// InitDefaultDB used for cli
func InitDefaultDB() *gorm.DB {
	return openDB(viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.addr"),
		viper.GetString("db.name"))
}

// GetDefaultDB 获取默认数据库示例
func GetDefaultDB() *gorm.DB {
	return InitDefaultDB()
}

// InitDockerDB 初始化一个docker数据库
func InitDockerDB() *gorm.DB {
	return openDB(viper.GetString("docker_db.username"),
		viper.GetString("docker_db.password"),
		viper.GetString("docker_db.addr"),
		viper.GetString("docker_db.name"))
}

// GetDockerDB 获取docker数据库
func GetDockerDB() *gorm.DB {
	return InitDockerDB()
}

// Init 初始化数据库
func (db *Database) Init() {
	DB = &Database{
		Default: GetDefaultDB(),
		Docker:  GetDockerDB(),
	}
}

// Close 关闭数据库链接
func (db *Database) Close() {
	err := DB.Default.Close()
	if err != nil {
		log.Warnf("[model] close default db err: %+v", err)
	}
	err = DB.Docker.Close()
	if err != nil {
		log.Warnf("[model] close docker db err: %+v", err)
	}
}
