package log

import (
	"errors"
	"fmt"

	"github.com/1024casts/snake/pkg/conf"
)

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
	errInvalidLoggerInstance = errors.New("log: invalid logger instance")
)

// Config is the struct for logger information
type Config struct {
	Name             string `yaml:"name"`
	Writers          string `yaml:"writers"`
	LoggerLevel      string `yaml:"logger_level"`
	LoggerFile       string `yaml:"logger_file"`
	LoggerWarnFile   string `yaml:"logger_warn_file"`
	LoggerErrorFile  string `yaml:"logger_error_file"`
	LogFormatText    bool   `yaml:"log_format_text"`
	LogRollingPolicy string `yaml:"log_rolling_policy"`
	LogRotateDate    int    `yaml:"log_rotate_date"`
	LogRotateSize    int    `yaml:"log_rotate_size"`
	LogBackupCount   uint   `yaml:"log_backup_count"`
}

// InitLog init log
func InitLog(cfg *conf.Config) Logger {
	logger, err := newZapLogger(cfg)
	if err != nil {
		fmt.Printf("InitWithConfig err: %v", err)
	}
	log = logger
	return logger
}

// Logger is our contract for the logger
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	WithFields(keyValues Fields) Logger
}

func GetLogger() Logger {
	return log
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
// output more field, eg:
// 		contextLogger := log.WithFields(log.Fields{"key1": "value1"})
// 		contextLogger.Info("print multi field")
// or more sample to use:
// 	    log.WithFields(log.Fields{"key1": "value1"}).Info("this is a test log")
// 	    log.WithFields(log.Fields{"key1": "value1"}).Infof("this is a test log, user_id: %d", userID)
func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}
