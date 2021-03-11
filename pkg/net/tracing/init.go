package tracing

import (
	"fmt"
	"io"

	"github.com/1024casts/snake/pkg/conf"

	"github.com/1024casts/snake/pkg/log"

	"github.com/uber/jaeger-lib/metrics"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// Init returns a new instance of Jaeger Tracer.
func Init(c *conf.Config, metricsFactory metrics.Factory) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: c.Jaeger.ServiceName,

		// "const" sampler is a binary sampling strategy: 0=never sample, 1=always sample.
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},

		// Log the emitted spans to stdout.
		Reporter: &config.ReporterConfig{
			LogSpans:           c.Jaeger.LogSpans,
			LocalAgentHostPort: c.Jaeger.Host,
		},
	}

	jaegerLogger := jaegerLoggerAdapter{log.GetLogger()}

	tracer, closer, err := cfg.NewTracer(
		//config.Logger(jaeger.StdLogger),
		config.Logger(jaegerLogger),
		config.Metrics(metricsFactory),
		config.ZipkinSharedRPCSpan(true),
	)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

type jaegerLoggerAdapter struct {
	logger log.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}
