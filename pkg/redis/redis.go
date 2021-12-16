package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-eagle/eagle/pkg/config"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
)

// RedisClient redis 客户端
var RedisClient *redis.Client

const (
	// ErrRedisNotFound not exist in redis
	ErrRedisNotFound = redis.Nil
	// DefaultRedisName default redis name
	DefaultRedisName = "default"
)

// RedisManager define a redis manager
//nolint
type RedisManager struct {
	clients map[string]*redis.Client
	*sync.RWMutex
}

// Init init a default redis instance
func Init() {
	clientManager := NewRedisManager()
	rdb, err := clientManager.GetClient(DefaultRedisName)
	if err != nil {
		panic(fmt.Sprintf("init redis err: %s", err.Error()))
	}
	RedisClient = rdb
}

// NewRedisManager create a redis manager
func NewRedisManager() *RedisManager {
	return &RedisManager{
		clients: make(map[string]*redis.Client),
		RWMutex: &sync.RWMutex{},
	}
}

// GetClient get a redis instance
func (r *RedisManager) GetClient(name string) (*redis.Client, error) {
	// get client from map
	r.RLock()
	if client, ok := r.clients[name]; ok {
		r.RUnlock()
		return client, nil
	}
	r.RUnlock()

	c, err := loadConf(name)
	if err != nil {
		panic(fmt.Sprintf("load redis conf err: %v", err))
	}

	// create a redis client
	r.Lock()
	defer r.Unlock()
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           c.DB,
		MinIdleConns: c.MinIdleConn,
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		PoolSize:     c.PoolSize,
		PoolTimeout:  c.PoolTimeout,
	})

	// check redis if is ok
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	// hook tracing (using open telemetry)
	if c.EnableTrace {
		rdb.AddHook(redisotel.NewTracingHook())
	}
	r.clients[name] = rdb

	return rdb, nil
}

// loadConf load redis config
func loadConf(name string) (ret *Config, err error) {
	v, err := config.LoadWithType("redis", "yaml")
	if err != nil {
		return nil, err
	}

	var c Config
	err = v.UnmarshalKey(name, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// InitTestRedis 实例化一个可以用于单元测试的redis
func InitTestRedis() {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	// 打开下面命令可以测试链接关闭的情况
	// defer mr.Close()

	RedisClient = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	fmt.Println("mini redis addr:", mr.Addr())
}
