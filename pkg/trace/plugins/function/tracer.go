package function

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	// TraceName function trance name
	TraceName = "github.com/go-eagle/eagle/trace/plugins/function"
	// PluginName plugin category name
	PluginName = "function"
)

var tracer trace.Tracer

func init() {
	tracer = otel.GetTracerProvider().Tracer(TraceName, trace.WithInstrumentationVersion(contrib.SemVersion()))
}

// StartFromContext create a new context
func StartFromContext(ctx context.Context) (context.Context, trace.Span) {
	var spanAttrs []trace.SpanStartOption

	fileName, lineNo, funcName := getCallerInfo(2)
	caller := fmt.Sprintf("FuncName: %s, file: %s, line: %d", funcName, fileName, lineNo)
	if caller != "" {
		spanAttrs = append(spanAttrs, trace.WithAttributes(attribute.String("caller", caller)))
	}

	spanName := fmt.Sprintf("%s: %s", PluginName, funcName)
	return tracer.Start(ctx, spanName, spanAttrs...)
}

func getCallerInfo(skip int) (string, int, string) {
	var (
		fileName string
		lineNo   int
		funcName string
	)
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		return fileName, lineNo, funcName
	}

	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(file)

	return fileName, lineNo, funcName
}
