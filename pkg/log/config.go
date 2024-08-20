package log

import "time"

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
	FlushInterval     time.Duration // default is 30s, recommend is dev or test is 1s, prod is 1m
}
