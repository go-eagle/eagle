package config

import (
	"errors"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf config

// config conf struct
type config struct {
	vp *viper.Viper
	// configDir conf root dir
	configDir string
	// configType conf file type, eg: yaml, json, toml, default yaml
	configType string
}

type Option func(*config)

// WithConfigDir config root dir
func WithConfigDir(cfgDir string) Option {
	return func(c *config) {
		c.configDir = cfgDir
	}
}

// WithFileType config dir
func WithFileType(typ string) Option {
	return func(c *config) {
		c.configType = typ
	}
}

// New create a config instance
func New(opts ...Option) *config {
	c := config{
		vp:         viper.New(),
		configType: "yaml",
	}
	for _, opt := range opts {
		opt(&c)
	}

	Conf = c

	return &c
}

// Load load file
func (c *config) Load(cfgName string) error {
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
	return nil
}

// Scan load data to val
func (c *config) Scan(val interface{}) error {
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
