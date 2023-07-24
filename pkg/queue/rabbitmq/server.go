package rabbitmq

import (
	"context"
	"sync"

	"github.com/go-eagle/eagle/pkg/log"
)

type SubscribeOptionMap map[string]Handler

type Server struct {
	subscriberOpts SubscribeOptionMap
	mu             sync.RWMutex
	started        bool
	consumer       *Consumer
}

func NewServer(cfg *Config) *Server {
	c, err := NewConsumer(cfg, log.GetLogger())
	if err != nil {
		panic(err)
	}
	srv := &Server{
		consumer:       c,
		subscriberOpts: SubscribeOptionMap{},
	}

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	if s.started {
		return nil
	}

	for _, h := range s.subscriberOpts {
		err := s.doConsume(ctx, h)
		if err != nil {
			return err
		}
	}

	s.started = true

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Info("[rabbitmq] server stopping...")
	s.started = false
	err := s.consumer.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) RegisterSubscriber(ctx context.Context, queueName string, h Handler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.subscriberOpts[queueName] = h

	return nil
}

func (s *Server) doConsume(ctx context.Context, h Handler) error {
	go func() {
		_ = s.consumer.Consume(ctx, h)
	}()

	return nil
}
