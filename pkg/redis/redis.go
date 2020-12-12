package redis

import (
	"fmt"

	"github.com/1024casts/snake/pkg/conf"
	"github.com/1024casts/snake/pkg/log"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
)

// RedisClient redis 客户端
var RedisClient *redis.Client

// Nil redis 返回为空
const Nil = redis.Nil

// Init 实例化一个redis client
func Init(cfg *conf.Config) *redis.Client {
	c := cfg.Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           c.Db,
		MinIdleConns: c.MinIdleConn,
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		PoolSize:     c.PoolSize,
		PoolTimeout:  c.PoolTimeout,
	})

	fmt.Println("redis addr:", cfg.Redis.Addr)

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
