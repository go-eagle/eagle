package http

// Option is a function that sets some option on the client.
type Option func(c *config)

// Options control behavior of the client.
type config struct {
	ClientTyp string
}

func WithClientType(clientType string) Option {
	return func(cfg *config) {
		cfg.ClientTyp = clientType
	}
}
