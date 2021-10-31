package app

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/registry"
	"github.com/go-eagle/eagle/pkg/transport"
)

// App global app
type App struct {
	opts     options
	ctx      context.Context
	cancel   func()
	mu       sync.Mutex
	instance *registry.ServiceInstance
}

// New create a app globally
func New(opts ...Option) *App {
	o := options{
		ctx:    context.Background(),
		logger: log.GetLogger(),
		// don not catch SIGKILL signal, need to waiting for kill self by other.
		sigs:            []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registryTimeout: 10 * time.Second,
	}
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(&o)
	}

	ctx, cancel := context.WithCancel(o.ctx)
	return &App{
		opts:   o,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Run start app
func (a *App) Run() error {
	// build service instance
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(a.ctx)

	// start server
	wg := sync.WaitGroup{}
	for _, srv := range a.opts.servers {
		srv := srv
		eg.Go(func() error {
			// wait for stop signal
			<-ctx.Done()
			return srv.Stop(ctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(ctx)
		})
	}

	// register service
	if a.opts.registry != nil {
		c, cancel := context.WithTimeout(a.opts.ctx, a.opts.registryTimeout)
		defer cancel()
		if err := a.opts.registry.Register(c, instance); err != nil {
			return err
		}
		a.mu.Lock()
		a.instance = instance
		a.mu.Unlock()
	}

	// watch signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, a.opts.sigs...)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case s := <-quit:
				a.opts.logger.Infof("receive a quit signal: %s", s.String())
				err := a.Stop()
				if err != nil {
					a.opts.logger.Infof("failed to stop app, err: %s", err.Error())
					return err
				}
			}
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

// Stop stops the application gracefully.
func (a *App) Stop() error {
	// deregister instance
	a.mu.Lock()
	instance := a.instance
	a.mu.Unlock()
	if a.opts.registry != nil && instance != nil {
		ctx, cancel := context.WithTimeout(a.opts.ctx, a.opts.registryTimeout)
		defer cancel()
		if err := a.opts.registry.Deregister(ctx, instance); err != nil {
			return err
		}
	}

	// cancel app
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

func (a *App) buildInstance() (*registry.ServiceInstance, error) {
	// register instance by withEndpoint
	endpoints := make([]string, 0)
	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}
	// auto register instance
	if len(endpoints) == 0 {
		for _, srv := range a.opts.servers {
			if r, ok := srv.(transport.Endpoint); ok {
				e, err := r.Endpoint()
				if err != nil {
					return nil, err
				}
				endpoints = append(endpoints, e.String())
			}
		}
	}
	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Version:   a.opts.version,
		Metadata:  a.opts.metadata,
		Endpoints: endpoints,
	}, nil
}
