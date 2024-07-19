package log

// Config  log config
type Config struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
	ServiceName       string // service name
	Filename          string
	Writers           string
	LoggerDir         string
	LogFormatText     bool
	LogRollingPolicy  string
	LogBackupCount    uint
}
