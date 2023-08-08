package main

import (
	"sync"
	"time"

	"github.com/go-eagle/eagle/pkg/config"
	"github.com/hibiken/asynq"
)

var (
	client *asynq.Client
	once   sync.Once
)

type Config struct {
	Queue struct {
		Addr         string
		Password     string
		DB           int
		MinIdleConn  int
		DialTimeout  time.Duration
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		PoolSize     int
		PoolTimeout  time.Duration
		Concurrency  int //并发数
	} `json:"redis"`
}

func GetClient() *asynq.Client {
	once.Do(func() {
		//c := config.New("config", config.WithEnv("local"))
		c := config.New(".")
		var cfg Config
		if err := c.Load("redis", &cfg); err != nil {
			panic(err)
		}
		client = asynq.NewClient(asynq.RedisClientOpt{
			Addr:         cfg.Queue.Addr,
			Password:     cfg.Queue.Password,
			DB:           cfg.Queue.DB,
			DialTimeout:  cfg.Queue.DialTimeout,
			ReadTimeout:  cfg.Queue.ReadTimeout,
			WriteTimeout: cfg.Queue.WriteTimeout,
			PoolSize:     cfg.Queue.PoolSize,
		})
	})
	return client
}
