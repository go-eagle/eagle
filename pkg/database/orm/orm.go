package orm

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// MySQL driver.
	"gorm.io/driver/mysql"
	// GORM MySQL
	"gorm.io/gorm"
)

// Config mysql config
type Config struct {
	Name            string
	Addr            string
	UserName        string
	Password        string
	ShowLog         bool
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifeTime time.Duration
}

// NewMySQL 链接数据库，生成数据库实例
func NewMySQL(c *Config) (db *gorm.DB) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		c.UserName,
		c.Password,
		c.Addr,
		c.Name,
		true,
		//"Asia/Shanghai"),
		"Local")

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panicf("open mysql failed. database name: %s, err: %+v", c.Name, err)
	}
	// set for db connection
	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(c.ConnMaxLifeTime)

	db, err = gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		log.Panicf("database connection failed. database name: %s, err: %+v", c.Name, err)
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	return db
}
