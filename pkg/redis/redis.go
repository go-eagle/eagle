package redis

import (
	"github.com/go-redis/redis"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

var Client *redis.Client

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	_, err := Client.Ping().Result()
	if err != nil {
		log.Errorf(err, "[redis] redis ping err: %+v")
		panic(err)
	}
}
