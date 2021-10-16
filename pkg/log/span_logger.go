// Package log span logger for trace
// reference: https://github.com/jaegertracing/jaeger/tree/master/examples/hotrod/pkg/log
package log

import (
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type spanLogger struct {
	logger     *zap.Logger
	span       trace.Span
	spanFields []zapcore.Field
}

func (sl spanLogger) Info(args ...interface{}) {
	msg := fmt.Sprint(args...)
	var fields []zap.Field
	sl.logToSpan("info", msg)
	sl.logger.Info(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Infof(format string, args ...interface{}) {
	msg := fmt.Sprint(format, args)
	var fields []zap.Field
	sl.logToSpan("Infof", msg)
	sl.logger.Info(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Warn(args ...interface{}) {
	msg := fmt.Sprint(args...)
	var fields []zap.Field
	sl.logToSpan("warn", msg)
	sl.logger.Warn(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Warnf(format string, args ...interface{}) {
	msg := fmt.Sprint(format, args)
	var fields []zap.Field
	sl.logToSpan("Warnf", msg)
	sl.span.RecordError(errors.New(msg))
	sl.logger.Warn(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Error(args ...interface{}) {
	msg := fmt.Sprint(args...)
	var fields []zap.Field
	sl.logToSpan("error", msg)
	sl.span.RecordError(errors.New(msg))
	sl.logger.Error(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprint(format, args)
	var fields []zap.Field
	sl.logToSpan("Errorf", msg)
	sl.span.RecordError(errors.New(msg))
	sl.logger.Error(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) WithFields(keyValues Fields) Logger {
	panic("implement me")
}

func (sl spanLogger) logToSpan(level string, msg string) {
	sl.span.SetAttributes(
		attribute.String("event", level),
		attribute.String("message", msg),
	)
}
