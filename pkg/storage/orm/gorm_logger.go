package orm

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/go-eagle/eagle/pkg/log"
)

// GormLogger is a custom logger for Gorm that uses the eagle log package.
type GormLogger struct {
	log                       log.Logger
	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	Colorful                  bool
}

// NewGormLogger creates a new GormLogger.
func NewGormLogger(l log.Logger, conf *Config) logger.Interface {
	logLevelMap := map[string]logger.LogLevel{
		"info":   logger.Info,
		"warn":   logger.Warn,
		"error":  logger.Error,
		"silent": logger.Silent,
	}
	logLevel, ok := logLevelMap[conf.LogLevel]
	if !ok {
		logLevel = logger.Info // if not found, use the default value
	}

	slowThreshold := 200 * time.Millisecond // Default value
	if conf.SlowThreshold != 0 {
		slowThreshold = conf.SlowThreshold
	}

	gl := &GormLogger{
		log:                       l,
		LogLevel:                  logLevel,
		SlowThreshold:             slowThreshold,
		IgnoreRecordNotFoundError: conf.IgnoreRecordNotFoundError,
		Colorful:                  conf.Colorful,
	}
	return gl
}

// LogMode sets the log level.
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level

	return &newLogger
}

// Info logs an info message.
func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.log.Infof(msg, args...)
	}
}

// Warn logs a warning message.
func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.log.Warnf(msg, args...)
	}
}

// Error logs an error message.
func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.log.Errorf(msg, args...)
	}
}

// Trace logs a SQL query.
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := log.Fields{
		"elapsed": elapsed,
		"sql":     sql,
		"rows":    rows,
	}

	switch {
	// log error
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		fields["err"] = err
		l.log.WithFields(fields).Error("gorm record not found")

	// log slow query
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= logger.Warn:
		fields["slow_threshold"] = l.SlowThreshold
		l.log.WithFields(fields).Warn("gorm slow query")

	// log all queries
	case l.LogLevel >= logger.Info:
		l.log.WithFields(fields).Debug("gorm trace")
	}
}
