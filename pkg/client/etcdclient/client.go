package etcdclient

import (
	"fmt"
	"time"

	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"

	"github.com/go-eagle/eagle/pkg/log"
)

// Config etcd config
type Config struct {
	Endpoints        []string      `json:"endpoints"`
	BasicAuth        bool          `json:"basicAuth"`
	UserName         string        `json:"userName"`
	Password         string        `json:"-"`
	ConnectTimeout   time.Duration `json:"connectTimeout"` // 连接超时时间
	Secure           bool          `json:"secure"`
	AutoSyncInterval time.Duration `json:"autoAsyncInterval"` // 自动同步member list的间隔
	TTL              int           // 单位：s
	logger           log.Logger
}

// Client ...
type Client struct {
	*clientv3.Client
	config *Config
}

func New(config *Config) (*Client, error) {
	return newClient(config)
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
		return nil, fmt.Errorf("client etcd endpoints empty, empty endpoints")
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
		// config.logger.Panic("client etcd start panic", xlog.FieldMod(ecode.ModClientETCD), xlog.FieldErrKind(ecode.ErrKindAny), xlog.FieldErr(err), xlog.FieldValueAny(config))
		return nil, fmt.Errorf("client etcd start failed: %v", err)
	}

	cc := &Client{
		Client: client,
		config: config,
	}

	config.logger.Info("dial etcd server")
	return cc, nil
}
