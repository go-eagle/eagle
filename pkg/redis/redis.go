package redis

import (
	"context"
	"fmt"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	apmgoredis "github.com/opentracing-contrib/goredis"

	"github.com/1024casts/snake/pkg/conf"
	"github.com/1024casts/snake/pkg/log"
)

// RedisClient redis 客户端
var RedisClient redis.UniversalClient

// Nil redis 返回为空
const Nil = redis.Nil

// Init 实例化一个redis client
func Init(cfg *conf.Config) redis.UniversalClient {
	c := cfg.Redis
	RedisClient = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        []string{c.Addr},
		Password:     c.Password,
		DB:           c.Db,
		MinIdleConns: c.MinIdleConn,
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		PoolSize:     c.PoolSize,
		PoolTimeout:  c.PoolTimeout,
	})

	// TODO: 使用每一次请求的context
	RedisClient = apmgoredis.Wrap(RedisClient).WithContext(context.Background())

	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Panicf("[redis] redis ping err: %+v", err)
	}
	return RedisClient
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
