package app

import (
	"context"
	"os"
	"time"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/transport"
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

// WithVersion with a version
func WithVersion(version string) Option {
	return func(o *options) {
		o.version = version
	}
}

// WithContext with a context
func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// WithSignal with some system signal
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

// WithServer with a server , http or grpc
func WithServer(srv ...transport.Server) Option {
	return func(o *options) {
		o.servers = srv
	}
}
