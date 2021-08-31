package trace

// Option is a function that sets some option on the client.
type Option func(c *Options)

// Options control behavior of the client.
type Options struct {
	SamplingRatio float64
}

func applyOptions(options ...Option) Options {
	opts := Options{
		1,
	}
	for _, option := range options {
		option(&opts)
	}

	return opts
}
