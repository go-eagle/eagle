package log

import "errors"

// A global variable so that log functions can be directly accessed
var log Logger

// Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

const (
	// InstanceZapLogger zap logger
	InstanceZapLogger int = iota
	// here add other logger
)

var (
	errInvalidLoggerInstance = errors.New("invalid logger instance")
)

// Logger is our contract for the logger
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	WithFields(keyValues Fields) Logger
}

// Config is the struct for logger information
type Config struct {
	Writers         string `yaml:"writers"`
	LoggerLevel     string `yaml:"logger_level"`
	LoggerFile      string `yaml:"logger_file"`
	LoggerWarnFile  string `yaml:"logger_warn_file"`
	LoggerErrorFile string `yaml:"logger_error_file"`
	LogFormatText   bool   `yaml:"log_format_text"`
	RollingPolicy   string `yaml:"rollingPolicy"`
	LogRotateDate   int    `yaml:"log_rotate_date"`
	LogRotateSize   int    `yaml:"log_rotate_size"`
	LogBackupCount  int    `yaml:"log_backup_count"`
}

// NewLogger returns an instance of logger
func NewLogger(cfg *Config, loggerInstance int) error {
	switch loggerInstance {
	case InstanceZapLogger:
		logger, err := newZapLogger(cfg)
		if err != nil {
			return err
		}
		log = logger
		return nil
	default:
		return errInvalidLoggerInstance
	}
}

// Debug logger
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Info logger
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warn logger
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Error logger
func Error(args ...interface{}) {
	log.Error(args...)
}

// Fatal logger
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Debugf logger
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Infof logger
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warnf logger
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Errorf logger
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatalf logger
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Panicf logger
func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

// WithFields logger
func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}
