package conf

import (
	"log"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/1024casts/snake/pkg/database/orm"
)

var (
	// Conf app global config
	Conf *Config
)

// Init init config
func Init(configPath string) (*Config, error) {
	cfgFile, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}
	cfg, err := ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	WatchConfig(cfgFile)

	Conf = cfg

	return cfg, nil
}

// LoadConfig load config file from given path
func LoadConfig(confPath string) (*viper.Viper, error) {
	v := viper.New()
	if confPath != "" {
		v.SetConfigFile(confPath) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		v.AddConfigPath("config") // 如果没有指定配置文件，则解析默认的配置文件
		v.SetConfigName("config")
	}
	v.SetConfigType("yaml")     // 设置配置文件格式为YAML
	v.AutomaticEnv()            // 读取匹配的环境变量
	viper.SetEnvPrefix("snake") // 读取环境变量的前缀为 snake
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

// 监控配置文件变化并热加载程序
func WatchConfig(v *viper.Viper) {
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
	})
}

// Config global config
// include common and biz config
type Config struct {
	// common
	App     AppConfig
	Logger  Logger
	MySQL   orm.Config
	Redis   RedisConfig
	Cache   CacheConfig
	Email   EmailConfig
	Web     WebConfig
	Cookie  CookieConfig
	QiNiu   QiNiuConfig
	Metrics Metrics
	Jaeger  Jaeger
	MongoDB MongoDB

	// here can add biz conf

}

// AppConfig app config
type AppConfig struct {
	Name              string
	Version           string
	Mode              string
	Port              string
	PprofPort         string
	URL               string
	JwtSecret         string
	JwtTimeout        int
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
}

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
	Name              string
	Writers           string
	LoggerFile        string
	LoggerWarnFile    string
	LoggerErrorFile   string
	LogFormatText     bool
	LogRollingPolicy  string
	LogRotateDate     int
	LogRotateSize     int
	LogBackupCount    uint
}

// RedisConfig redis config
type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	MinIdleConn  int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
	PoolTimeout  time.Duration
}

// CacheConfig define cache config struct
type CacheConfig struct {
	Driver string
	Prefix string
}

// EmailConfig email config
type EmailConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	Name      string
	Address   string
	ReplyTo   string
	KeepAlive int
}

// WebConfig web config
type WebConfig struct {
	Name   string
	Domain string
	Secret string
	Static string
}

// CookieConfig cookie config
type CookieConfig struct {
	Name     string
	MaxAge   int
	Secure   bool
	HttpOnly bool
	Domain   string
	Secret   string
}

// QiNiuConfig qiniu config
type QiNiuConfig struct {
	AccessKey   string
	SecretKey   string
	CdnURL      string
	SignatureID string
	TemplateID  string
}

// Metrics config
type Metrics struct {
	URL         string
	ServiceName string
}

// Jaeger config
type Jaeger struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

// MongoDB config
type MongoDB struct {
	URI      string
	User     string
	Password string
	DB       string
}
