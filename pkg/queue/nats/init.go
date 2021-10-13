package nats

var (
	Queue *Producer
)

type Config struct {
	Addr string `mapstructure:"name"`
}

func Init(cfg *Config) {
	Queue = NewProducer(cfg.Addr)
}
