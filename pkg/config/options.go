package config

// Option config option
type Option func(*config)

// WithConfigDir config root dir
func WithConfigDir(cfgDir string) Option {
	return func(c *config) {
		c.configDir = cfgDir
	}
}

// WithFileType config dir
func WithFileType(typ string) Option {
	return func(c *config) {
		c.configType = typ
	}
}
