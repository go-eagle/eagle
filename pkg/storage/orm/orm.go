package orm

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"

	"github.com/go-eagle/eagle/pkg/config"
)

const (
	// DriverMySQL mysql driver
	DriverMySQL = "mysql"
	// DriverPostgres postgresSQL driver
	DriverPostgres = "postgres"
	// DriverClickhouse
	DriverClickhouse = "clickhouse"

	// DefaultDatabase default db name
	DefaultDatabase = "default"
)

var (
	// DBMap store database instance
	DBMap = make(map[string]*gorm.DB)
	// DBLock database locker
	DBLock sync.Mutex
	// logWriter log writer
	LogWriter logger.Writer
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
	Timeout         string // connect timeout
	ReadTimeout     string
	WriteTimeout    string
	ConnMaxLifeTime time.Duration
	SlowThreshold   time.Duration // 慢查询时长，默认500ms
	EnableTrace     bool
}

// New create a or multi database client
func New(names ...string) error {
	if len(names) == 0 {
		return fmt.Errorf("no set databasename")
	}

	clientManager := NewManager()
	for _, name := range names {
		_, err := clientManager.GetInstance(name)
		if err != nil {
			return fmt.Errorf("init database name: %+v, err: %+v", name, err)
		}
	}

	return nil
}

// Manager define a manager
type Manager struct {
	instances map[string]*gorm.DB
	*sync.RWMutex
}

// NewManager create a database manager
func NewManager() *Manager {
	return &Manager{
		instances: make(map[string]*gorm.DB),
		RWMutex:   &sync.RWMutex{},
	}
}

// GetDB get a database
func GetDB(name string) (*gorm.DB, error) {
	DBLock.Lock()
	defer DBLock.Unlock()

	db, ok := DBMap[name]
	if !ok {
		db, err := NewManager().GetInstance(name)
		if err != nil {
			return nil, err
		}
		return db, nil
	}

	return db, nil
}

// GetInstance return a database client
func (m *Manager) GetInstance(name string) (*gorm.DB, error) {
	// get client from map
	m.RLock()
	if ins, ok := m.instances[name]; ok {
		m.RUnlock()
		return ins, nil
	}
	m.RUnlock()

	c, err := LoadConf(name)
	if err != nil {
		return nil, fmt.Errorf("load database conf err: %+v", err)
	}

	// create a database client
	m.Lock()
	defer m.Unlock()

	instance := NewInstance(c)
	m.instances[name] = instance
	DBMap[name] = instance

	return instance, nil
}

// NewInstance connect to database and create a db instance
func NewInstance(c *Config) (db *gorm.DB) {
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
	case DriverClickhouse:
		db, err = gorm.Open(clickhouse.Open(dsn), gormConfig(c))
	default:
		db, err = gorm.Open(mysql.Open(dsn), gormConfig(c))
	}
	if err != nil {
		log.Panicf("open db failed. driver: %s, database name: %s, err: %+v", c.Driver, c.Name, err)
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

	// set trace
	if c.EnableTrace {
		err = db.Use(tracing.NewPlugin())
		if err != nil {
			log.Panicf("using gorm opentelemetry, err: %+v", err)
		}
	}

	return db
}

// LoadConf load database config
func LoadConf(name string) (ret *Config, err error) {
	v, err := config.LoadWithType("database", "yaml")
	if err != nil {
		return nil, err
	}

	var c Config
	err = v.UnmarshalKey(name, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// getDSN return dsn string
func getDSN(c *Config) string {
	// default mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=%s&readTimeout=%s&writeTimeout=%s",
		c.UserName,
		c.Password,
		c.Addr,
		c.Name,
		c.Timeout,
		c.ReadTimeout,
		c.WriteTimeout,
	)

	if c.Driver == DriverPostgres {
		dsn = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable&connect_timeout=%s&statement_timeout=%s",
			c.UserName,
			c.Password,
			c.Addr,
			c.Name,
			c.Timeout,
			c.ReadTimeout,
		)
	}

	if c.Driver == DriverClickhouse {
		dsn = fmt.Sprintf("clickhouse://%s:%s@%s/%s?dial_timeout=%s&read_timeout=%s",
			c.UserName,
			c.Password,
			c.Addr,
			c.Name,
			c.Timeout,
			c.ReadTimeout,
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
		var writer logger.Writer
		//将标准输出作为Writer
		writer = log.New(os.Stdout, "\r\n", log.LstdFlags)
		// use custom logger
		if LogWriter != nil {
			writer = LogWriter
		}

		// new logger with writer
		config.Logger = logger.New(
			writer,
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
