package tracing

import (
	"fmt"
	"io"
	"time"

	"github.com/1024casts/snake/pkg/log"

	"github.com/uber/jaeger-lib/metrics"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// Jaeger config
type Config struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

// Init returns a new instance of Jaeger Tracer.
func Init(serviceName, agentHostPort string, metricsFactory metrics.Factory) (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: serviceName,

		// "const" sampler is a binary sampling strategy: 0=never sample, 1=always sample.
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},

		// Log the emitted spans to stdout.
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  agentHostPort,
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
		return nil, nil, err
	}
	opentracing.SetGlobalTracer(tracer)

	return tracer, closer, err
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
