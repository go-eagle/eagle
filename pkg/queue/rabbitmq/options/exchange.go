package options

type ExchangeOptions struct {
	Name       string                 `yaml:"name"`
	Kind       string                 `yaml:"kind"`
	Durable    bool                   `yaml:"durable"`
	AutoDelete bool                   `yaml:"auto-delete"`
	Internal   bool                   `yaml:"internal"`
	NoWait     bool                   `yaml:"no-wait"`
	Args       map[string]interface{} `yaml:"args"`
}
