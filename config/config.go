package config

import (
	"strings"
	"sync"

	"github.com/1024casts/snake/pkg/errno"

	"github.com/fsnotify/fsnotify"
	"github.com/lexkong/log"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	c   configure
	mux sync.Mutex
)

type configure struct {
	DB dbConfigure `mapstructure:"db"`
}

func (c *configure) Validate() error {
	if err := c.DB.Validate(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

type Config struct {
	Name string
}

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}

	// 初始化日志包
	c.initLog()

	// 监控配置文件变化并热加载程序
	c.watchConfig()

	return nil
}

func (cfg *Config) initConfig() error {
	if cfg.Name != "" {
		viper.SetConfigFile(cfg.Name) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		viper.AddConfigPath("conf") // 如果没有指定配置文件，则解析默认的配置文件
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml") // 设置配置文件格式为YAML
	viper.AutomaticEnv()        // 读取匹配的环境变量
	viper.SetEnvPrefix("snake") // 读取环境变量的前缀为goapi
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return errors.WithStack(err)
	}

	mux.Lock()
	defer mux.Unlock()
	if err := viper.Unmarshal(&c); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (cfg *Config) initLog() {
	passLagerCfg := log.PassLagerCfg{
		Writers:        viper.GetString("log.writers"),
		LoggerLevel:    viper.GetString("log.logger_level"),
		LoggerFile:     viper.GetString("log.logger_file"),
		LogFormatText:  viper.GetBool("log.log_format_text"),
		RollingPolicy:  viper.GetString("log.rollingPolicy"),
		LogRotateDate:  viper.GetInt("log.log_rotate_date"),
		LogRotateSize:  viper.GetInt("log.log_rotate_size"),
		LogBackupCount: viper.GetInt("log.log_backup_count"),
	}

	log.InitWithConfig(&passLagerCfg)
}

// 监控配置文件变化并热加载程序
func (cfg *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}

type dbConfigure struct {
	Host     string `mapstructure:"host"`
	Name     string `mapstructure:"name"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Port     uint   `mapstructure:"port"`
	Verbose  bool   `mapstructure:"verbose"`
}

func (c *dbConfigure) GetHost() string {
	return c.Host
}

func (c *dbConfigure) GetPort() uint {
	return c.Port
}

func (c *dbConfigure) GetUsername() string {
	return c.Username
}

func (c *dbConfigure) GetPassword() string {
	return c.Password
}

func (c *dbConfigure) GetDatabaseName() string {
	return c.Name
}

func (c *dbConfigure) GetVerbose() bool {
	return c.Verbose
}

func (c *dbConfigure) Validate() error {
	if c.Host == "" {
		return errors.Wrap(errno.ErrParam, "db.host")
	}
	if c.Port == 0 {
		return errors.Wrap(errno.ErrParam, "db.password")
	}
	if c.Username == "" {
		return errors.Wrap(errno.ErrParam, "db.username")
	}
	if c.Password == "" {
		return errors.Wrap(errno.ErrParam, "db.password")
	}
	if c.Name == "" {
		return errors.Wrap(errno.ErrParam, "db.name")
	}
	return nil
}

func GetServerConfig() configure {
	return c
}
