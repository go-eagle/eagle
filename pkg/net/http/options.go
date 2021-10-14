package http

import "time"

// Option is a function that sets some option on the client.
type Option func(c *option)

// Options control behavior of the client.
type option struct {
	ClientTyp string
	header    map[string][]string
	// timeout of per request
	timeout time.Duration
}

func WithClientType(clientType string) Option {
	return func(cfg *option) {
		cfg.ClientTyp = clientType
	}
}
