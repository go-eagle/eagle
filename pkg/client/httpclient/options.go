package httpclient

import "time"

// Option is a function that sets some option on the client.
type Option func(c *options)

// Options control behavior of the client.
type options struct {
	header map[string][]string
	// timeout of per request
	timeout time.Duration
}

func defaultOptions() *options {
	return &options{
		header:  make(map[string][]string),
		timeout: DefaultTimeout,
	}
}

// WithTimeout with a timeout for per request
func WithTimeout(duration time.Duration) Option {
	return func(cfg *options) {
		cfg.timeout = duration
	}
}
