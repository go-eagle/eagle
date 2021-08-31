package trace

import (
	"errors"
	"strings"

	jaegerprop "go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// InitTracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func InitTracerProvider(serviceName, endpoint string, options ...Option) (*tracesdk.TracerProvider, error) {
	var endpointOption jaeger.EndpointOption
	if serviceName == "" {
		return nil, errors.New("no service name provided")
	}
	if strings.HasPrefix(endpoint, "http") {
		endpointOption = jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint))
	} else {
		endpointOption = jaeger.WithAgentEndpoint(jaeger.WithAgentHost(endpoint))
	}

	// Create the Jaeger exporter
	exporter, err := jaeger.New(endpointOption)
	if err != nil {
		return nil, err
	}

	opts := applyOptions(options...)
	tp := tracesdk.NewTracerProvider(
		// set sample
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(opts.SamplingRatio)),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(jaegerprop.Jaeger{})

	return tp, nil
}
