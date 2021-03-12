package log

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Factory is the default logging wrapper that can create
// logger instances either for a given Context or context-less.
type Factory struct {
	zapLogger *zap.Logger
	logger    Logger
}

// NewFactory creates a new Factory.
func NewFactory(zapLogger *zap.Logger, logger Logger) Factory {
	return Factory{
		zapLogger: zapLogger,
		logger:    logger,
	}
}

// Bg creates a context-unaware logger.
func (b Factory) Bg() Logger {
	return b.logger
}

// For returns a context-aware Logger. If the context
// contains an OpenTracing span, all logging calls are also
// echo-ed into the span.
func (b Factory) For(ctx context.Context) Logger {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		logger := spanLogger{span: span, logger: b.zapLogger}

		if jaegerCtx, ok := span.Context().(jaeger.SpanContext); ok {
			logger.spanFields = []zapcore.Field{
				zap.String("trace_id", jaegerCtx.TraceID().String()),
				zap.String("span_id", jaegerCtx.SpanID().String()),
			}
		}

		return logger
	}
	return b.Bg()
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (b Factory) With(fields ...zapcore.Field) Factory {
	return Factory{zapLogger: b.zapLogger.With(fields...)}
}
