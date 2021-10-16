package nats

var (
	// Queue nats queue
	Queue *Producer
)

// Config nats config
type Config struct {
	Addr string `mapstructure:"name"`
}

// Init create producer
func Init(cfg *Config) {
	Queue = NewProducer(cfg.Addr)
}
