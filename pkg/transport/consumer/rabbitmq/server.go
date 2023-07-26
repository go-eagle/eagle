package rabbitmq

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/queue/rabbitmq"
	"github.com/go-eagle/eagle/pkg/queue/rabbitmq/options"
)

type SubscribeMap map[string]rabbitmq.Handler

type Server struct {
	subscribers SubscribeMap
	mu          sync.RWMutex
	opts        []options.ConsumerOption
	stop        chan struct{}
}

func NewServer(opts ...options.ConsumerOption) (*Server, error) {
	conf := rabbitmq.GetConfig()
	if len(conf) == 0 {
		return nil, fmt.Errorf("rabbitmq config is empty")
	}
	srv := &Server{
		opts:        opts,
		subscribers: make(SubscribeMap),
	}

	for name, _ := range conf {
		srv.subscribers[name] = nil
	}

	return srv, nil
}

func (s *Server) Start(ctx context.Context) error {
	var wg sync.WaitGroup
	for name, h := range s.subscribers {
		wg.Add(1)
		go func(taskName string, handler rabbitmq.Handler) {
			defer func() {
				wg.Done()
				log.Infof("[rabbitmq] task %s is stopped", taskName)
			}()

			done := make(chan struct{})
			consumer, err := rabbitmq.NewConsumer(rabbitmq.GetConfig()[taskName], log.GetLogger())
			if err != nil {
				log.Errorf("[rabbitmq] new consumer %s error: %v", taskName, err)
				return
			}
			go func() {
				defer close(done)
				err := consumer.Consume(ctx, handler, s.opts...)
				if err != nil {
					log.Errorf("[rabbitmq] start consume %s error: %v", taskName, err)
				}
			}()

			log.Infof("[rabbitmq] task %s is started successfully", taskName)

			select {
			case <-ctx.Done():
				log.Infof("[rabbitmq] task %s is stopping via context", taskName)
				return
			case <-s.stop:
				log.Infof("[rabbitmq] task %s is stopping via stop channel", taskName)
				_ = consumer.Close()
				return
			case <-done:
				return
			}
		}(name, h)
	}
	wg.Wait()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Info("[rabbitmq] server stopping...")
	close(s.stop)
	return nil
}

func (s *Server) RegisterSubscriber(ctx context.Context, queueName string, h rabbitmq.Handler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.subscribers[queueName] = h

	return nil
}
