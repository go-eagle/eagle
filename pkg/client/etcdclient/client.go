package etcdclient

import (
	"fmt"
	"time"

	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"

	"github.com/go-eagle/eagle/pkg/config"
	"github.com/go-eagle/eagle/pkg/log"
)

// Config etcd config
type Config struct {
	Endpoints        []string
	BasicAuth        bool
	UserName         string
	Password         string
	ConnectTimeout   time.Duration // 连接超时时间
	Secure           bool
	AutoSyncInterval time.Duration // 自动同步member list的间隔
	TTL              int           // 单位：s
	logger           log.Logger
}

// Client ...
type Client struct {
	*clientv3.Client
	config *Config
}

func New() (*Client, error) {
	cfg, err := loadConf()
	if err != nil {
		panic(fmt.Sprintf("[etcd] load etcd conf err: %v", err))
	}
	return newClient(cfg)
}

// New ...
func newClient(config *Config) (*Client, error) {
	conf := clientv3.Config{
		Endpoints:            config.Endpoints,
		DialTimeout:          config.ConnectTimeout,
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithUnaryInterceptor(grpcprom.UnaryClientInterceptor),
			grpc.WithStreamInterceptor(grpcprom.StreamClientInterceptor),
		},
		AutoSyncInterval: config.AutoSyncInterval,
	}

	config.logger = log.GetLogger()

	if config.Endpoints == nil {
		return nil, fmt.Errorf("[etcd]  client etcd endpoints empty, empty endpoints")
	}

	if !config.Secure {
		conf.DialOptions = append(conf.DialOptions, grpc.WithInsecure())
	}

	if config.BasicAuth {
		conf.Username = config.UserName
		conf.Password = config.Password
	}

	client, err := clientv3.New(conf)

	if err != nil {
		return nil, fmt.Errorf("[etcd] client etcd start failed: %v", err)
	}

	cc := &Client{
		Client: client,
		config: config,
	}

	config.logger.Info("dial etcd server")
	return cc, nil
}

// loadConf load register config
func loadConf() (ret *Config, err error) {
	var cfg Config
	v, err := config.LoadWithType("registry", config.FileTypeYaml)
	if err != nil {
		return nil, err
	}

	err = v.UnmarshalKey("etcd", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
