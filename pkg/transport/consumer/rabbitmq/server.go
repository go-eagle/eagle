package rabbitmq

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/queue/rabbitmq"
	"github.com/go-eagle/eagle/pkg/queue/rabbitmq/options"
)

var (
	// Srv is a global rabbitmq server
	Srv *Server
)

type HandlerMap map[string]rabbitmq.Handler

type Server struct {
	handlers HandlerMap
	mu       sync.RWMutex
	opts     []options.ConsumerOption
	stop     chan struct{}
}

func NewServer(opts ...options.ConsumerOption) *Server {
	conf := rabbitmq.GetConfig()
	if len(conf) == 0 {
		panic(fmt.Errorf("rabbitmq config is empty"))
	}
	srv := &Server{
		opts:     opts,
		handlers: make(HandlerMap),
		stop:     make(chan struct{}, 1),
	}

	for name, _ := range conf {
		srv.handlers[name] = nil
	}

	Srv = srv

	return srv
}

func GetServer() *Server {
	return Srv
}

func (s *Server) Start(ctx context.Context) error {
	var wg sync.WaitGroup
	for name, _ := range s.handlers {
		wg.Add(1)

		// get handler and check it
		h, err := s.GetRegisterHandler(name)
		if err != nil {
			log.Errorf("[rabbitmq] get handler %s error: %v", name, err)
			return err
		}

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
				log.Infof("[rabbitmq] task %s is stopping by cancel", taskName)
				_ = consumer.Close()
				select {
				case <-done:
					return
				}
			case <-s.stop:
				log.Infof("[rabbitmq] receiving stop signal, task %s is stopping", taskName)
				_ = consumer.Close()
			}
		}(name, h)
	}

	wg.Wait()

	log.Infof("[rabbitmq] all tasks are stopped successfully")

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Info("[rabbitmq] server stopping...")
	s.stop <- struct{}{}
	return nil
}

func (s *Server) RegisterHandler(queueName string, h rabbitmq.Handler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.handlers[queueName] = h

	return nil
}

func (s *Server) GetRegisterHandler(queueName string) (rabbitmq.Handler, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	handler, ok := s.handlers[queueName]
	if !ok {
		return nil, fmt.Errorf("handler %s not found", queueName)
	}

	return handler, nil
}
