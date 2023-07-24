package options

type ConsumerOption func(*ConsumerOptions)

type ConsumerOptions struct {
	Concurrency int

	QOSPrefetchCount int
	QOSPrefetchSize  int
	QOSGlobal        bool

	ConsumerName      string
	ConsumerAutoAck   bool
	ConsumerExclusive bool
	ConsumerNoWait    bool
	ConsumerNoLocal   bool
	ConsumerArgs      map[string]interface{}
	ConsumerQueue     string
}

func NewConsumerOptions(opts ...ConsumerOption) *ConsumerOptions {
	// set default value
	options := ConsumerOptions{
		Concurrency:      1,
		QOSPrefetchCount: 5,
	}

	// apply options
	for _, opt := range opts {
		opt(&options)
	}
	return &options
}

// WithConsumerOptionConcurrency sets concurrency
func WithConsumerOptionConcurrency(concurrency int) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.Concurrency = concurrency
	}
}

// WithConsumerOptionQOSPrefetch sets QOSPrefetchCount
func WithConsumerOptionQOSPrefetch(prefetchCount int) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.QOSPrefetchCount = prefetchCount
	}
}

// WithConsumerOptionQOSPrefetchSize sets QOSPrefetchSize
func WithConsumerOptionQOSPrefetchSize(prefetchSize int) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.QOSPrefetchSize = prefetchSize
	}
}

// WithConsumerOptionQOSGlobal sets QOSGlobal
func WithConsumerOptionQOSGlobal(global bool) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.QOSGlobal = global
	}
}

// WithConsumerOptionConsumerName sets ConsumerName
func WithConsumerOptionConsumerName(consumerName string) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.ConsumerName = consumerName
	}
}

// WithConsumerOptionConsumerAutoAck sets ConsumerAutoAck
func WithConsumerOptionConsumerAutoAck(autoAck bool) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.ConsumerAutoAck = autoAck
	}
}

// WithConsumerOptionConsumerExclusive sets ConsumerExclusive
func WithConsumerOptionConsumerExclusive(exclusive bool) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.ConsumerExclusive = exclusive
	}
}

// WithConsumerOptionConsumerNoWait sets ConsumerNoWait
func WithConsumerOptionConsumerNoWait(noWait bool) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.ConsumerNoWait = noWait
	}
}

// WithConsumerOptionConsumerNoLocal sets ConsumerNoLocal
func WithConsumerOptionConsumerNoLocal(noLocal bool) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.ConsumerNoLocal = noLocal
	}
}

// WithConsumerOptionConsumerArgs sets ConsumerArgs
func WithConsumerOptionConsumerArgs(args map[string]interface{}) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.ConsumerArgs = args
	}
}

// WithConsumerOptionConsumerQueue sets ConsumerQueue
func WithConsumerOptionConsumerQueue(queue string) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.ConsumerQueue = queue
	}
}
