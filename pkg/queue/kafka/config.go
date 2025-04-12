package kafka

import (
	"sync"

	cfg "github.com/go-eagle/eagle/pkg/config"
)

var (
	loadOnce  sync.Once
	closeOnce sync.Once
	conf      map[string]*Conf
)

type Conf struct {
	Version      string   `yaml:"Version"`
	RequiredAcks int      `yaml:"RequiredAcks"`
	Topic        string   `yaml:"Topic"`
	ConsumeTopic []string `yaml:"VonsumeTopic"`
	Brokers      []string `yaml:"Brokers"`
	GroupID      string   `yaml:"GroupID"`
	// partitioner type，optional: "random", "roundrobin", "hash"
	Partitioner string `yaml:"Partitioner"`
}

// loadConf load config
func loadConf() (ret map[string]*Conf, err error) {
	v, err := cfg.LoadWithType("kafka", "yaml")
	if err != nil {
		return nil, err
	}

	c := make(map[string]*Conf, 0)
	err = v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	conf = c

	return c, nil
}

func GetConfig() map[string]*Conf {
	return conf
}

func Load() {
	loadOnce.Do(func() {
		conf, err := loadConf()
		if err != nil {
			panic(err)
		}

		DefaultManager = NewManager(conf)
	})
}

func Close() {
	closeOnce.Do(func() {
		_ = DefaultManager.Close()
	})
}
