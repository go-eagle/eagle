package options

type QueueOptions struct {
	Name       string                 `yaml:"name"`
	Durable    bool                   `yaml:"durable"`
	AutoDelete bool                   `yaml:"auto-delete"`
	Exclusive  bool                   `yaml:"exclusive"`
	NoWait     bool                   `yaml:"no-wait"`
	Args       map[string]interface{} `yaml:"args"`
}
