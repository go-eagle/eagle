package breaker

import (
	"sync"
	"time"

	xtime "github.com/go-eagle/eagle/pkg/time"
)

// Breaker is a CircuitBreaker pattern.
// FIXME on int32 atomic.LoadInt32(&b.on) == _switchOn
type Breaker interface {
	Allow() error
	MarkSuccess()
	MarkFailed()
}

const (
	// StateOpen when circuit breaker open, request not allowed, after sleep
	// some duration, allow one single request for testing the health, if ok
	// then state reset to closed, if not continue the step.
	StateOpen int32 = iota
	// StateClosed when circuit breaker closed, request allowed, the breaker
	// calc the succeed ratio, if request num greater request setting and
	// ratio lower than the setting ratio, then reset state to open.
	StateClosed
	// StateHalfopen when circuit breaker open, after slepp some duration, allow
	// one request, but not state closed.
	StateHalfopen

	//_switchOn int32 = iota
	// _switchOff
)

var (
	_mu   sync.RWMutex
	_conf = &Config{
		Window:  xtime.Duration(3 * time.Second),
		Bucket:  10,
		Request: 100,

		// Percentage of failures must be lower than 33.33%
		K: 1.5,

		// Pattern: "",
	}
	_group = NewGroup(_conf)
)

// Init init global breaker config, also can reload config after first time call.
func Init(conf *Config) {
	if conf == nil {
		return
	}
	_mu.Lock()
	_conf = conf
	_mu.Unlock()
}

// newBreaker new a breaker.
func newBreaker(c *Config) (b Breaker) {
	// factory
	return newSRE(c)
}

// Go runs your function while tracking the breaker state of default group.
func Go(name string, run, fallback func() error) error {
	breaker := _group.Get(name)
	if err := breaker.Allow(); err != nil {
		return fallback()
	}
	return run()
}
