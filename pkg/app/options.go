package app

import (
	"context"
	"os"
	"time"

	"github.com/1024casts/snake/pkg/log"
	"github.com/1024casts/snake/pkg/transport"
)

// Option is func for application
type Option func(o *options)

// options is an application options
type options struct {
	id      string
	name    string
	version string

	sigs []os.Signal
	ctx  context.Context

	logger           log.Logger
	registrarTimeout time.Duration
	servers          []transport.Server
}

// WithID with app id
func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

// WithName .
func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

// WithVersion .
func WithVersion(version string) Option {
	return func(o *options) {
		o.version = version
	}
}

// WithSignal .
func WithSignal(sigs ...os.Signal) Option {
	return func(o *options) {
		o.sigs = sigs
	}
}

// WithLogger .
func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func Server(srv ...transport.Server) Option {
	return func(o *options) {
		o.servers = srv
	}
}
