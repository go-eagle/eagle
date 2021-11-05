package config

// Option config option
type Option func(*config)

// WithConfigDir config root dir
func WithConfigDir(cfgDir string) Option {
	return func(c *config) {
		c.configDir = cfgDir
	}
}

// WithFileType config file type
func WithFileType(fileType string) Option {
	return func(c *config) {
		c.configType = fileType
	}
}

// WithEnv env var
func WithEnv(name string) Option {
	return func(c *config) {
		c.env = name
	}
}
