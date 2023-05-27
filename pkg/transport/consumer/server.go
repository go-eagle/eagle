package consumer

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/hibiken/asynq"
)

type Server struct {
	clientOpt asynq.RedisClientOpt

	//  async server
	srv *asynq.Server
	mux *asynq.ServeMux

	// async schedule
	sche *asynq.Scheduler
}

func NewServer(redisOpt asynq.RedisClientOpt, asyncCfg asynq.Config) *Server {
	srv := &Server{
		srv: asynq.NewServer(redisOpt, asyncCfg),
		mux: asynq.NewServeMux(),
		sche: asynq.NewScheduler(
			redisOpt,
			&asynq.SchedulerOpts{Location: time.Local},
		),
	}

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	err := s.srv.Run(s.mux)
	if err != nil {
		return errors.Wrapf(err, "failed to run async server: %v")
	}
	err = s.sche.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to run async Scheduler server: %v")
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.srv.Shutdown()
	s.sche.Shutdown()
	return nil
}

// RegisterTask register task
func (s *Server) RegisterTask(schedule string, task *asynq.Task) (entryID string, err error) {
	return s.sche.Register(schedule, task)
}

// RegisterHandle register handler
func (s *Server) RegisterHandle(pattern string, handler func(context.Context, *asynq.Task) error) {
	s.mux.HandleFunc(pattern, handler)
}
