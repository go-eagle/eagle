package options

type BindOptions struct {
	RoutingKey string                 `yaml:"routing-key"`
	NoWait     bool                   `yaml:"no-wait"`
	Args       map[string]interface{} `yaml:"args"`
}
