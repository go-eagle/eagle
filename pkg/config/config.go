package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	// Conf global conf var
	Conf *config
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

// CommonConfig app config.
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

// ServerConfig server config.
type ServerConfig struct {
	Network      string
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// config conf struct.
type config struct {
	// env environment var
	env string
	// configDir conf root dir
	configDir string
	// configType conf file type, eg: yaml, json, toml, default yaml
	configType string
	val        map[string]*viper.Viper
	mu         sync.Mutex
}

// New create a config instance.
func New(opts ...Option) *config {
	c := config{
		configType: "yaml",
		val:        make(map[string]*viper.Viper),
	}
	for _, opt := range opts {
		opt(&c)
	}

	Conf = &c

	return &c
}

// Scan scan data to struct.
func (c *config) Scan(filename string, val interface{}) error {
	v, err := c.LoadByType(filename, c.configType)
	if err != nil {
		return err
	}
	err = v.Unmarshal(&val)
	if err != nil {
		return err
	}
	return nil
}

// LoadByType load conf by file type.
func (c *config) LoadByType(filename string, cfgType string) (v *viper.Viper, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.val[filename]
	if ok {
		return v, nil
	}

	v, err = c.load(filename, cfgType)
	if err != nil {
		return nil, err
	}
	c.val[filename] = v
	return v, nil
}

// Scan load file.
func (c *config) load(filename string, cfgType string) (*viper.Viper, error) {
	// application parameters take precedence over environment variables
	env := GetEnvString("APP_ENV", "")
	path := filepath.Join(c.configDir, env)
	if c.env != "" {
		path = filepath.Join(c.configDir, c.env)
	}

	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName(filename)
	v.SetConfigType(c.configType)
	if cfgType != "" {
		v.SetConfigType(cfgType)
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
	})

	return v, nil
}

// GetEnvString get value from env.
func GetEnvString(key string, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}
