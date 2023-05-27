package cron

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
)

type Server struct {
	clientOpt asynq.RedisClientOpt
	sche      *asynq.Scheduler

	tasks map[string]*asynq.Task
}

func NewServer(redisOpt asynq.RedisClientOpt) *Server {
	srv := &Server{
		sche: asynq.NewScheduler(
			redisOpt,
			&asynq.SchedulerOpts{Location: time.Local},
		),
	}

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	err := s.sche.Run()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.sche.Shutdown()
	return nil
}

func (s *Server) RegisterTask(schedule string, task *asynq.Task) (entryID string, err error) {
	return s.sche.Register(schedule, task)
}
