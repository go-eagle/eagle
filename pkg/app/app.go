package app

import (
	"context"
	"os"
	"syscall"

	"github.com/google/uuid"
)

type App struct {
	opts   options
	ctx    context.Context
	cancel func()
}

func New(opts ...Option) *App {
	options := options{
		logger: nil,
		sigs:   []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
	if id, err := uuid.NewUUID(); err == nil {
		options.id = id.String()
	}
	for _, o := range opts {
		o(&options)
	}

	ctx, cancel := context.WithCancel(options.ctx)
	return &App{
		opts:   options,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (a *App) Run() error {

	return nil
}

func (a *App) Stop() error {

	return nil
}
