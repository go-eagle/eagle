package consulclient

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/go-eagle/eagle/pkg/config"
)

// Config consul config
type Config struct {
	Addr       string
	Scheme     string
	Datacenter string
	WaitTime   time.Duration
	Namespace  string
}

func New() (*api.Client, error) {
	cfg, err := loadConf()
	if err != nil {
		panic(fmt.Sprintf("[etcd] load consul conf err: %v", err))
	}
	return newClient(cfg)
}

func newClient(cfg *Config) (*api.Client, error) {
	consulClient, err := api.NewClient(&api.Config{
		Address:    cfg.Addr,
		Scheme:     cfg.Scheme,
		Datacenter: cfg.Datacenter,
		WaitTime:   cfg.WaitTime,
		Namespace:  cfg.Namespace,
	})
	if err != nil {
		log.Fatal(err)
	}
	return consulClient, nil
}

// loadConf load register config
func loadConf() (ret *Config, err error) {
	var cfg Config
	v, err := config.LoadWithType("registry", config.FileTypeYaml)
	if err != nil {
		return nil, err
	}

	err = v.UnmarshalKey("consul", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
