package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-eagle/eagle/pkg/config"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
)

// RedisClient redis 客户端
var RedisClient *redis.Client

// ErrRedisNotFound not exist in redis
const ErrRedisNotFound = redis.Nil

// Config redis config
type Config struct {
	Addr         string
	Password     string
	DB           int
	MinIdleConn  int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
	PoolTimeout  time.Duration
	// tracing switch
	EnableTrace bool
}

// Init 实例化一个redis client
func Init() *redis.Client {
	c, err := loadConf()
	if err != nil {
		panic(fmt.Sprintf("load redis conf err: %v", err))
	}

	RedisClient = redis.NewClient(&redis.Options{
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

	_, err = RedisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	// hook tracing (using open telemetry)
	if c.EnableTrace {
		RedisClient.AddHook(redisotel.NewTracingHook())
	}

	return RedisClient
}

// loadConf load redis config
func loadConf() (ret *Config, err error) {
	var cfg Config
	if err := config.Load("redis", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
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
