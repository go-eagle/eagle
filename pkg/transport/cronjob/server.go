package cronjob

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
)

type Server struct {
	clientOpt asynq.RedisClientOpt
	sche      *asynq.Scheduler

	jobs map[string]*asynq.Task
}

func NewServer() *Server {
	srv := &Server{
		sche: asynq.NewScheduler(
			asynq.RedisClientOpt{Addr: "127.0.0.1:6379"},
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

func (s *Server) RegisterJob(schedule string, job *asynq.Task) (entryID string, err error) {
	return s.sche.Register(schedule, job)
}
