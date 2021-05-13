package trace

import (
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/1024casts/snake/pkg/log"
)

// TraceAgent is an abstraction for trace agent (Zipkin, Jaeger, Elastic).
type TraceAgent interface {
	New(componentName string) (opentracing.Tracer, io.Closer, error)
}

type Trace struct {
	ServiceName string

	tracer opentracing.Tracer
	closer io.Closer
}

func Init(serviceName string, traceAgent TraceAgent) (*Trace, error) {
	trace := &Trace{
		ServiceName: serviceName,
	}

	if traceAgent == nil {
		return nil, fmt.Errorf("not set trace agent")
	}

	var err error
	trace.tracer, trace.closer, err = traceAgent.New(serviceName)
	if err != nil {
		return nil, err
	}
	return trace, nil
}

// StartSpan delegates to opentracing.Tracer.
func (t *Trace) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return t.tracer.StartSpan(operationName, opts...)
}

// Inject delegates to opentracing.Tracer.
func (t *Trace) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return t.tracer.Inject(sm, format, carrier)
}

// Extract delegates to opentracing.Tracer.
func (t *Trace) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return t.tracer.Extract(format, carrier)
}

// IsEnabled determines if Tracer was successfully activated.
func (t *Trace) IsEnabled() bool {
	return t != nil && t.tracer != nil
}

// Close trace
func (t *Trace) Close() {
	if t.closer != nil {
		err := t.closer.Close()
		if err != nil {
			log.Error("close trace error, %v", err)
		}
	}
}

// LogRequest used to create span tags from the request.
func LogRequest(span opentracing.Span, r *http.Request) {
	if span != nil && r != nil && r.URL != nil {
		ext.HTTPMethod.Set(span, r.Method)
		ext.HTTPUrl.Set(span, r.URL.String())
		span.SetTag("http.host", r.Host)
	}
}

// LogResponseCode used to log response code in span.
func LogResponseCode(span opentracing.Span, code int) {
	if span != nil {
		ext.HTTPStatusCode.Set(span, uint16(code))
	}
}

// LogEventf logs an event to the span in the request context.
func LogEventf(span opentracing.Span, format string, args ...interface{}) {
	if span != nil {
		span.LogKV("event", fmt.Sprintf(format, args...))
	}
}

// SetError flags the span associated with this request as in error.
func SetError(span opentracing.Span) {
	if span != nil {
		ext.Error.Set(span, true)
	}
}

// SetErrorWithEvent flags the span associated with this request as in error and log an event.
func SetErrorWithEvent(span opentracing.Span, format string, args ...interface{}) {
	SetError(span)
	LogEventf(span, format, args...)
}
