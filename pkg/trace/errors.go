package trace

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SetSpanError record error to tracing system
func SetSpanError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)

	if span.IsRecording() {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}
