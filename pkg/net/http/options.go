package http

import "time"

// Option is a function that sets some option on the client.
type Option func(c *option)

// Options control behavior of the client.
type option struct {
	header map[string][]string
	// timeout of per request
	timeout time.Duration
}

func defaultOptions() *option {
	return &option{
		header:  make(map[string][]string),
		timeout: DefaultTimeout,
	}
}

// WithTimeout with a timeout for per request
func WithTimeout(duration time.Duration) Option {
	return func(cfg *option) {
		cfg.timeout = duration
	}
}
