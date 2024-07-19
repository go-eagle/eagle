package log

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/go-eagle/eagle/pkg/config"
)

// log is A global variable so that log functions can be directly accessed
var log Logger
var logger Logger
var zl *zap.Logger

// Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

// Logger is a contract for the logger
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	WithFields(keyValues Fields) Logger
}

// loadConf load logger config
func loadConf() (ret *Config, err error) {
	var cfg Config
	if err := config.Load("logger", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Init init log
func Init(opts ...Option) Logger {
	var err error
	cfg, err := loadConf()
	if err != nil {
		panic(fmt.Sprintf("load logger conf err: %v", err))
	}

	// new zap logger
	zl, err = newZapLogger(cfg, opts...)
	if err != nil {
		_ = fmt.Errorf("init newZapLogger err: %v", err)
	}
	_ = zl

	// log 用于支持模块级的方法调用，所以要比其他 Logger 多跳一层
	log, err = newLoggerWithCallerSkip(cfg, 1, opts...)
	if err != nil {
		_ = fmt.Errorf("init newLogger err: %v", err)
	}

	logger, err = newLogger(cfg, opts...)
	if err != nil {
		_ = fmt.Errorf("init logger err: %v", err)
	}

	return log
}

// GetLogger return a log
func GetLogger() Logger {
	return logger
}

// GetZapLogger return raw zap logger
func GetZapLogger() *zap.Logger {
	return zl
}

// WithContext is a logger that can log msg and log span for trace
func WithContext(ctx context.Context) Logger {
	//return zap logger

	if span := trace.SpanFromContext(ctx); span != nil {
		logger := spanLogger{span: span, logger: zl}

		spanCtx := span.SpanContext()
		logger.spanFields = []zapcore.Field{
			zap.String("trace_id", spanCtx.TraceID().String()),
			zap.String("span_id", spanCtx.SpanID().String()),
		}

		return logger
	}
	return GetLogger()
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

// WithFields logger
// output more field, eg:
//
//	contextLogger := log.WithFields(log.Fields{"key1": "value1"})
//	contextLogger.Info("print multi field")
//
// or more sample to use:
//
//	log.WithFields(log.Fields{"key1": "value1"}).Info("this is a test log")
//	log.WithFields(log.Fields{"key1": "value1"}).Infof("this is a test log, user_id: %d", userID)
func WithFields(keyValues Fields) Logger {
	return GetLogger().WithFields(keyValues)
}
