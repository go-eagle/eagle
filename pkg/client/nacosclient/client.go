package nacosclient

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/go-eagle/eagle/pkg/config"
)

// Config consul config
type Config struct {
	Addr        string
	Port        uint64
	NamespaceId string // default public
	TimeoutMs   uint64 // unit: default 10000ms
	LogDir      string
	CacheDir    string
	LogLevel    string // default value is info
}

func New() (naming_client.INamingClient, error) {
	cfg, err := loadConf()
	if err != nil {
		panic(fmt.Sprintf("[etcd] load consul conf err: %v", err))
	}
	return newClient(cfg)
}

func newClient(cfg *Config) (naming_client.INamingClient, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(cfg.Addr, cfg.Port),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         cfg.NamespaceId,
		TimeoutMs:           cfg.TimeoutMs,
		NotLoadCacheAtStart: true,
		LogDir:              cfg.LogDir,
		CacheDir:            cfg.CacheDir,
		LogLevel:            cfg.LogLevel,
	}

	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		return nil, err
	}

	return cli, nil
}

// loadConf load register config
func loadConf() (ret *Config, err error) {
	var cfg Config
	v, err := config.LoadWithType("registry", config.FileTypeYaml)
	if err != nil {
		return nil, err
	}

	err = v.UnmarshalKey("nacos", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
