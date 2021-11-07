package config

// Option config option
type Option func(*Config)

// WithFileType config file type
func WithFileType(fileType string) Option {
	return func(c *Config) {
		c.configType = fileType
	}
}

// WithEnv env var
func WithEnv(name string) Option {
	return func(c *Config) {
		c.env = name
	}
}
