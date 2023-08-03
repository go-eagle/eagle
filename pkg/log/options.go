package log

type Option func(*Config)

// WithFilename set log filename
func WithFilename(filename string) Option {
	return func(cfg *Config) {
		cfg.Name = filename
	}
}

// WithLogDir set log dir
func WithLogDir(dir string) Option {
	return func(cfg *Config) {
		cfg.LoggerDir = dir
	}
}
