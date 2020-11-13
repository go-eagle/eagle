package conf

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/1024casts/snake/pkg/log"
)

var (
	// Conf app global config
	Conf *Config
)

// Init init config
func Init(confPath string) error {
	err := initConfig(confPath)
	if err != nil {
		return err
	}
	return nil
}

// initConfig init config from conf file
func initConfig(confPath string) error {
	if confPath != "" {
		viper.SetConfigFile(confPath) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		viper.AddConfigPath("conf") // 如果没有指定配置文件，则解析默认的配置文件
		viper.SetConfigName("config.local")
	}
	viper.SetConfigType("yaml") // 设置配置文件格式为YAML
	viper.AutomaticEnv()        // 读取匹配的环境变量
	viper.SetEnvPrefix("snake") // 读取环境变量的前缀为 snake
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return errors.WithStack(err)
	}

	// parse to config struct
	err := viper.Unmarshal(&Conf)
	if err != nil {
		return err
	}

	watchConfig()

	return nil
}

// 监控配置文件变化并热加载程序
func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}

// Config global config
// include common and biz config
type Config struct {
	// common
	App   AppConfig
	Log   LogConfig
	MySQL MySQLConfig
	Redis RedisConfig
	Cache CacheConfig

	// here can add biz conf
}

// AppConfig app config
type AppConfig struct {
	Name      string
	RunMode   string
	Addr      string
	URL       string
	JwtSecret string
}

// LogConfig log config
type LogConfig struct {
	Writers          string
	LoggerLevel      string
	LoggerFile       string
	LoggerWarnFile   string
	LoggerErrorFile  string
	LogFormatText    bool
	LogRollingPolicy string
	LogRotateDate    int
	LogRotateSize    int
	LogBackupCount   int
}

// MySQLConfig mysql config
type MySQLConfig struct {
	Name            string
	Addr            string
	UserName        string
	Password        string
	ShowLog         bool
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifeTime int
}

// RedisConfig redis config
type RedisConfig struct {
	Addr         string
	Password     string
	Db           int
	minIdleConn  int
	DialTimeout  int
	ReadTimeout  int
	WriteTimeout int
	PoolSize     int
	poolTimeout  int
}

// CacheConfig define cache config struct
type CacheConfig struct {
	Driver string
	Prefix string
}

// InitLog init log
func InitLog() {
	config := log.Config{
		Writers:          viper.GetString("log.writers"),
		LoggerLevel:      viper.GetString("log.logger_level"),
		LoggerFile:       viper.GetString("log.logger_file"),
		LoggerWarnFile:   viper.GetString("log.logger_warn_file"),
		LoggerErrorFile:  viper.GetString("log.logger_error_file"),
		LogFormatText:    viper.GetBool("log.log_format_text"),
		LogRollingPolicy: viper.GetString("log.log_rolling_policy"),
		LogRotateDate:    viper.GetInt("log.log_rotate_date"),
		LogRotateSize:    viper.GetInt("log.log_rotate_size"),
		LogBackupCount:   viper.GetInt("log.log_backup_count"),
	}
	err := log.NewLogger(&config, log.InstanceZapLogger)
	if err != nil {
		fmt.Printf("InitWithConfig err: %v", err)
	}
}
