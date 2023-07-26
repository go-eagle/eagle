package rabbitmq

import (
	"sync"
	"time"

	"github.com/go-eagle/eagle/pkg/config"

	"github.com/go-eagle/eagle/pkg/queue/rabbitmq/options"
)

var (
	loadOnce  sync.Once
	closeOnce sync.Once
	conf      map[string]*Config
)

type Config struct {
	URI         string                     `yaml:"uri"`
	AutoDeclare bool                       `yaml:"auto-declare"`
	Timeout     time.Duration              `yaml:"timeout"`
	Connection  *options.ConnectionOptions `yaml:"connection"`
	Exchange    *options.ExchangeOptions   `yaml:"exchange"`
	Queue       *options.QueueOptions      `yaml:"queue"`
	Bind        *options.BindOptions       `yaml:"bind"`
}

// loadConf load config
func loadConf() (ret map[string]*Config, err error) {
	v, err := config.LoadWithType("rabbitmq", "yaml")
	if err != nil {
		return nil, err
	}

	c := make(map[string]*Config, 0)
	err = v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	conf = c

	return c, nil
}

func GetConfig() map[string]*Config {
	return conf
}

func Load() {
	loadOnce.Do(func() {
		conf, err := loadConf()
		if err != nil {
			panic(err)
		}

		for _, v := range conf {
			conn, err := options.NewConnectionOptions(options.WithURI(v.URI))
			if err != nil {
				panic(err)
			}
			v.Connection = conn
		}

		DefaultManager = NewManager(conf)
		if err != nil {
			panic(err)
		}
	})
}

func Close() {
	closeOnce.Do(func() {
		_ = DefaultManager.Close()
	})
}
