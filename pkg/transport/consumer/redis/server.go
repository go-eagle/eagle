package redis

import (
	"context"

	"github.com/pkg/errors"

	"github.com/hibiken/asynq"
)

const (
	// QueueCritical queue priority
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

// Server async server
type Server struct {
	clientOpt asynq.RedisClientOpt

	//  async server
	srv *asynq.Server
	mux *asynq.ServeMux
}

// NewServer new async server
func NewServer(redisOpt asynq.RedisClientOpt, asyncCfg asynq.Config) *Server {
	srv := &Server{
		srv: asynq.NewServer(redisOpt, asyncCfg),
		mux: asynq.NewServeMux(),
	}

	return srv
}

// Start async server
func (s *Server) Start(ctx context.Context) error {
	err := s.srv.Run(s.mux)
	if err != nil {
		return errors.Wrapf(err, "failed to run async server")
	}

	return nil
}

// Stop async server
func (s *Server) Stop(ctx context.Context) error {
	s.srv.Shutdown()
	return nil
}

// RegisterHandler register handler
func (s *Server) RegisterHandler(pattern string, handler func(context.Context, *asynq.Task) error) {
	s.mux.HandleFunc(pattern, handler)
}
