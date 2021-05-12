// span logger for tracing
// reference: https://github.com/jaegertracing/jaeger/tree/master/examples/hotrod/pkg/log
package log

import (
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	tag "github.com/opentracing/opentracing-go/ext"
	spanlog "github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type spanLogger struct {
	logger     *zap.Logger
	span       opentracing.Span
	spanFields []zapcore.Field
}

func (sl spanLogger) Debug(args ...interface{}) {
	panic("implement me")
}

func (sl spanLogger) Debugf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	var fields []zap.Field
	sl.logToSpan("info", msg)
	sl.logger.Debug(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Info(args ...interface{}) {
	msg := fmt.Sprint(args...)
	var fields []zap.Field
	sl.logToSpan("info", msg)
	sl.logger.Info(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Infof(format string, args ...interface{}) {
	panic("implement me")
}

func (sl spanLogger) Warn(args ...interface{}) {
	msg := fmt.Sprint(args...)
	var fields []zap.Field
	sl.logToSpan("error", msg)
	sl.logger.Warn(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Warnf(format string, args ...interface{}) {
	msg := fmt.Sprint(format, args)
	var fields []zap.Field
	sl.logToSpan("error", msg)
	sl.logger.Warn(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Error(args ...interface{}) {
	msg := fmt.Sprint(args...)
	var fields []zap.Field
	sl.logToSpan("error", msg)
	sl.logger.Error(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Errorf(format string, args ...interface{}) {
	panic("implement me")
}

func (sl spanLogger) Fatal(args ...interface{}) {
	msg := fmt.Sprint(args...)
	var fields []zap.Field
	sl.logToSpan("fatal", msg)
	tag.Error.Set(sl.span, true)
	sl.logger.Fatal(msg, append(sl.spanFields, fields...)...)
}

func (sl spanLogger) Fatalf(format string, args ...interface{}) {
	panic("implement me")
}

func (sl spanLogger) Panicf(format string, args ...interface{}) {
	panic("implement me")
}

func (sl spanLogger) WithFields(keyValues Fields) Logger {
	panic("implement me")
}

func (sl spanLogger) logToSpan(level string, msg string) {
	// TODO rather than always converting the fields, we could wrap them into a lazy logger
	fa := fieldAdapter(make([]spanlog.Field, 0, 2))
	fa = append(fa, spanlog.String("event", msg))
	fa = append(fa, spanlog.String("level", level))
	sl.span.LogFields(fa...)
}

type fieldAdapter []spanlog.Field

func (fa *fieldAdapter) AddBool(key string, value bool) {
	*fa = append(*fa, spanlog.Bool(key, value))
}

func (fa *fieldAdapter) AddFloat64(key string, value float64) {
	*fa = append(*fa, spanlog.Float64(key, value))
}

func (fa *fieldAdapter) AddFloat32(key string, value float32) {
	*fa = append(*fa, spanlog.Float64(key, float64(value)))
}

func (fa *fieldAdapter) AddInt(key string, value int) {
	*fa = append(*fa, spanlog.Int(key, value))
}

func (fa *fieldAdapter) AddInt64(key string, value int64) {
	*fa = append(*fa, spanlog.Int64(key, value))
}

func (fa *fieldAdapter) AddInt32(key string, value int32) {
	*fa = append(*fa, spanlog.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddInt16(key string, value int16) {
	*fa = append(*fa, spanlog.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddInt8(key string, value int8) {
	*fa = append(*fa, spanlog.Int64(key, int64(value)))
}

func (fa *fieldAdapter) AddUint(key string, value uint) {
	*fa = append(*fa, spanlog.Uint64(key, uint64(value)))
}

func (fa *fieldAdapter) AddUint64(key string, value uint64) {
	*fa = append(*fa, spanlog.Uint64(key, value))
}

func (fa *fieldAdapter) AddUint32(key string, value uint32) {
	*fa = append(*fa, spanlog.Uint64(key, uint64(value)))
}

func (fa *fieldAdapter) AddUint16(key string, value uint16) {
	*fa = append(*fa, spanlog.Uint64(key, uint64(value)))
}

func (fa *fieldAdapter) AddUint8(key string, value uint8) {
	*fa = append(*fa, spanlog.Uint64(key, uint64(value)))
}

func (fa *fieldAdapter) AddUintptr(key string, value uintptr)                        {}
func (fa *fieldAdapter) AddArray(key string, marshaler zapcore.ArrayMarshaler) error { return nil }
func (fa *fieldAdapter) AddComplex128(key string, value complex128)                  {}
func (fa *fieldAdapter) AddComplex64(key string, value complex64)                    {}
func (fa *fieldAdapter) AddObject(key string, value zapcore.ObjectMarshaler) error   { return nil }
func (fa *fieldAdapter) AddReflected(key string, value interface{}) error            { return nil }
func (fa *fieldAdapter) OpenNamespace(key string)                                    {}

func (fa *fieldAdapter) AddDuration(key string, value time.Duration) {
	// TODO inefficient
	*fa = append(*fa, spanlog.String(key, value.String()))
}

func (fa *fieldAdapter) AddTime(key string, value time.Time) {
	// TODO inefficient
	*fa = append(*fa, spanlog.String(key, value.String()))
}

func (fa *fieldAdapter) AddBinary(key string, value []byte) {
	*fa = append(*fa, spanlog.Object(key, value))
}

func (fa *fieldAdapter) AddByteString(key string, value []byte) {
	*fa = append(*fa, spanlog.Object(key, value))
}

func (fa *fieldAdapter) AddString(key, value string) {
	if key != "" && value != "" {
		*fa = append(*fa, spanlog.String(key, value))
	}
}
