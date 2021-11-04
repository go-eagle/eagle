package config

import (
	"errors"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	// Conf global conf var
	Conf Config
	// App global app var
	App AppConfig
)

// Config global config
// nolint
type AppConfig struct {
	CommonConfig
	HTTP ServerConfig
	GRPC ServerConfig
}

// CommonConfig app config
type CommonConfig struct {
	Name              string
	Version           string
	Mode              string
	PprofPort         string
	URL               string
	JwtSecret         string
	JwtTimeout        int
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
	EnableTrace       bool
}

// ServerConfig server config
type ServerConfig struct {
	Network      string
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Config define config interface
type Config interface {
	Load(cfgName string, val interface{}) error
}

// config conf struct
type config struct {
	vp *viper.Viper
	// configDir conf root dir
	configDir string
	// configType conf file type, eg: yaml, json, toml, default yaml
	configType string
}

// New create a config instance
func New(opts ...Option) Config {
	c := config{
		vp:         viper.New(),
		configType: "yaml",
	}
	for _, opt := range opts {
		opt(&c)
	}

	Conf = &c

	return &c
}

// Load load file
func (c *config) Load(cfgName string, val interface{}) error {
	c.vp.AddConfigPath(c.configDir)
	if cfgName != "" {
		c.vp.SetConfigName(cfgName)
	}
	c.vp.SetConfigType(c.configType)

	if err := c.vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.New("config file not found")
		}
		return err
	}

	err := c.vp.Unmarshal(&val)
	if err != nil {
		return err
	}

	c.watch(&val)

	return nil
}

// watch listen file change
func (c *config) watch(v interface{}) {
	go func() {
		c.vp.WatchConfig()
		c.vp.OnConfigChange(func(e fsnotify.Event) {
			_ = c.vp.Unmarshal(&v)
			log.Printf("Config file changed: %s", e.Name)
		})
	}()
}
