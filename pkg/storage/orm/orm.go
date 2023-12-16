package orm

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	otelgorm "github.com/1024casts/gorm-opentelemetry"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"

	// GORM MySQL
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	// DriverMySQL mysql driver
	DriverMySQL = "mysql"
	// DriverPostgres postgresSQL driver
	DriverPostgres = "postgres"
)

// Config database config
type Config struct {
	Driver          string
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

// New connect to database and create a db instance
func New(c *Config) (db *gorm.DB) {
	var (
		err   error
		sqlDB *sql.DB
	)
	dsn := getDSN(c)
	switch c.Driver {
	case DriverMySQL:
		db, err = gorm.Open(mysql.Open(dsn), gormConfig(c))
	case DriverPostgres:
		db, err = gorm.Open(postgres.Open(dsn), gormConfig(c))
	default:
		db, err = gorm.Open(mysql.Open(dsn), gormConfig(c))
	}
	if err != nil {
		log.Panicf("open mysql failed. driver: %s, database name: %s, err: %+v", c.Driver, c.Name, err)
	}

	sqlDB, err = db.DB()
	if err != nil {
		log.Panicf("database connection failed. database name: %s, err: %+v", c.Name, err)
	}
	// set for db connection
	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(c.ConnMaxLifeTime)

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

// getDSN return dsn string
func getDSN(c *Config) string {
	// default mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		c.UserName,
		c.Password,
		c.Addr,
		c.Name,
		true,
		//"Asia/Shanghai"),
		"Local")

	if c.Driver == DriverPostgres {
		dsn = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			c.UserName,
			c.Password,
			c.Addr,
			c.Name,
		)
	}

	return dsn
}

// gormConfig 根据配置决定是否开启日志
func gormConfig(c *Config) *gorm.Config {
	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true} // 禁止外键约束, 生产环境不建议使用外键约束
	// 打印所有SQL
	if c.ShowLog {
		config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}
	// 只打印慢查询
	if c.SlowThreshold > 0 {
		config.Logger = logger.New(
			//将标准输出作为Writer
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				//设定慢查询时间阈值
				SlowThreshold: c.SlowThreshold, // nolint: golint
				Colorful:      true,
				//设置日志级别，只有指定级别以上会输出慢查询日志
				LogLevel: logger.Warn,
			},
		)
	}
	return config
}
