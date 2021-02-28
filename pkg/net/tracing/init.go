package tracing

import (
	"fmt"
	"io"

	"github.com/uber/jaeger-lib/metrics"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// Init returns a new instance of Jaeger Tracer.
func Init(serviceName string, metricsFactory metrics.Factory) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: serviceName,

		// "const" sampler is a binary sampling strategy: 0=never sample, 1=always sample.
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},

		// Log the emitted spans to stdout.
		Reporter: &config.ReporterConfig{
			LogSpans: true,
			//LocalAgentHostPort:  "127.0.0.1:6381",
			//BufferFlushInterval: 100 * time.Millisecond,
			//CollectorEndpoint:   "http://127.0.0.1:14268/api/traces",   // for gorm
		},
	}
	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaeger.StdLogger),
		config.Metrics(metricsFactory),
		config.ZipkinSharedRPCSpan(true),
	)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}
