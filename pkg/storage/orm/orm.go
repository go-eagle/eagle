package orm

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"

	otelgorm "github.com/1024casts/gorm-opentelemetry"

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
	SlowThreshold   time.Duration // 慢查询时长，默认500ms
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

	var gormLogger logger.Interface
	// 打印所有SQL
	if c.ShowLog {
		gormLogger = logger.Default.LogMode(logger.Info)
	}
	// 只打印慢查询
	if c.SlowThreshold > 0 {
		gormLogger = logger.New(
			//将标准输出作为Writer
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				//设定慢查询时间阈值
				SlowThreshold: c.SlowThreshold,
				Colorful:      true,
				//设置日志级别，只有指定级别以上会输出慢查询日志
				LogLevel: logger.Warn,
			},
		)
	}

	db, err = gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Panicf("database connection failed. database name: %s, err: %+v", c.Name, err)
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	// Initialize otel plugin with options
	plugin := otelgorm.NewPlugin(
	// include any options here
	)

	// set trace
	err = db.Use(plugin)
	if err != nil {
		log.Panicf("using gorm opentelemetry, err: %+v", err)
	}

	return db
}
