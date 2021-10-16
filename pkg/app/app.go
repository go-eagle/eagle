package app

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/go-eagle/eagle/pkg/conf"
	"github.com/go-eagle/eagle/pkg/log"
)

// App global app
type App struct {
	c      *conf.Config
	opts   options
	ctx    context.Context
	cancel func()
	log    log.Logger
}

// New create a app globally
func New(c *conf.Config, opts ...Option) *App {
	defaultOptions := options{
		ctx:    context.Background(),
		logger: log.GetLogger(),
		// don not catch SIGKILL signal, need to waiting for kill self by other.
		sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout: 10 * time.Second,
	}
	if id, err := uuid.NewUUID(); err == nil {
		defaultOptions.id = id.String()
	}
	for _, o := range opts {
		o(&defaultOptions)
	}

	ctx, cancel := context.WithCancel(defaultOptions.ctx)
	return &App{
		c:      c,
		opts:   defaultOptions,
		ctx:    ctx,
		log:    log.GetLogger(),
		cancel: cancel,
	}
}

// Run start app
func (a *App) Run() error {
	a.log.Infof("app_name: %s, version: %s", a.opts.name, a.opts.version)
	eg, ctx := errgroup.WithContext(a.ctx)

	// start server
	for _, srv := range a.opts.servers {
		srv := srv
		eg.Go(func() error {
			// wait for stop signal
			<-ctx.Done()
			return srv.Stop(ctx)
		})
		eg.Go(func() error {
			return srv.Start(ctx)
		})
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
				a.log.Infof("receive a quit signal: %s", s.String())
				return a.Stop()
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
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}
