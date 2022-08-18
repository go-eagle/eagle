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

func NewServer(addr, exchangeName string) *Server {
	srv := &Server{
		consumer:       NewConsumer(addr, exchangeName, false),
		subscriberOpts: SubscribeOptionMap{},
	}

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	if s.started {
		return nil
	}

	err := s.consumer.Start()
	if err != nil {
		return err
	}

	log.Infof("[rabbitmq] server listening on: %s", s.consumer.addr)

	for k, h := range s.subscriberOpts {
		err := s.doConsume(ctx, k, h)
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
	s.consumer.Stop()
	return nil
}

func (s *Server) RegisterSubscriber(ctx context.Context, queueName string, h Handler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.subscriberOpts[queueName] = h

	return nil
}

func (s *Server) doConsume(ctx context.Context, queueName string, h Handler) error {
	go s.consumer.Consume(ctx, queueName, h)

	return nil
}
