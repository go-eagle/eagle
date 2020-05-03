package redis

import (
	"fmt"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"

	"github.com/1024casts/snake/pkg/log"
)

// Client redis 客户端
var Client *redis.Client

// Nil redis 返回为空
const Nil = redis.Nil

// Init 实例化一个redis client
func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"),
	})

	fmt.Println("redis addr:", viper.GetString("redis.addr"))

	_, err := Client.Ping().Result()
	if err != nil {
		log.Errorf("[redis] redis ping err: %+v", err)
		panic(err)
	}
}

// InitTestRedis 实例化一个可以用于单元测试的redis
func InitTestRedis() {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	// 打开下面命令可以测试链接关闭的情况
	// defer mr.Close()

	Client = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	fmt.Println("mini redis addr:", mr.Addr())
}
