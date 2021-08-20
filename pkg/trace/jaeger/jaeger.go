package jaeger

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	jaegercli "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/zipkin"
	jaegermet "github.com/uber/jaeger-lib/metrics"

	"github.com/go-eagle/eagle/pkg/log"
)

// Name sets the name of this tracer.
const Name = "jaeger"

// Config provides configuration settings for a jaeger tracer.
type Config struct {
	SamplingServerURL      string  // Set the sampling server url
	SamplingType           string  // Set the sampling type
	SamplingParam          float64 // Set the sampling parameter
	LocalAgentHostPort     string  // Set jaeger-agent's host:port that the reporter will used
	Gen128Bit              bool    // Generate 128 bit span IDs
	Propagation            string  // Which propagation format to use (jaeger/b3)
	TraceContextHeaderName string  // Set the header to use for the trace-id
	CollectorEndpoint      string  // Instructs reporter to send spans to jaeger-collector at this URL
	CollectorUser          string  // CollectorUser for basic http authentication when sending spans to jaeger-collector
	CollectorPassword      string  // CollectorPassword for basic http authentication when sending spans to jaeger-collector
}

// SetDefaults sets the default values.
func (c *Config) SetDefaults() {
	c.SamplingServerURL = "http://localhost:5778/sampling"
	c.SamplingType = "const"
	c.SamplingParam = 1.0
	c.LocalAgentHostPort = "127.0.0.1:6831"
	c.Propagation = "jaeger"
	c.Gen128Bit = true
	c.TraceContextHeaderName = jaegercli.TraceContextHeaderName
	c.CollectorEndpoint = ""
	c.CollectorUser = ""
	c.CollectorPassword = ""
}

// New sets up the tracer
func (c *Config) New(componentName string) (opentracing.Tracer, io.Closer, error) {
	reporter := &jaegercfg.ReporterConfig{
		LogSpans:           true,
		LocalAgentHostPort: c.LocalAgentHostPort,
		CollectorEndpoint:  c.CollectorEndpoint,
		User:               c.CollectorUser,
		Password:           c.CollectorPassword,
	}

	jcfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			SamplingServerURL: c.SamplingServerURL,
			Type:              c.SamplingType,
			Param:             c.SamplingParam,
		},
		Reporter: reporter,
		Headers: &jaegercli.HeadersConfig{
			TraceContextHeaderName: c.TraceContextHeaderName,
		},
	}

	jMetricsFactory := jaegermet.NullFactory

	opts := []jaegercfg.Option{
		jaegercfg.Metrics(jMetricsFactory),
		jaegercfg.Gen128Bit(c.Gen128Bit),
	}

	switch c.Propagation {
	case "b3":
		p := zipkin.NewZipkinB3HTTPHeaderPropagator()
		opts = append(opts,
			jaegercfg.Injector(opentracing.HTTPHeaders, p),
			jaegercfg.Extractor(opentracing.HTTPHeaders, p),
		)
	case "jaeger", "":
	default:
		return nil, nil, fmt.Errorf("unknown propagation format: %s", c.Propagation)
	}

	// Initialize tracer with a logger and a metrics factory
	closer, err := jcfg.InitGlobalTracer(
		componentName,
		opts...,
	)
	if err != nil {
		log.Errorf("Could not initialize jaeger tracer: %s", err.Error())
		return nil, nil, err
	}
	return opentracing.GlobalTracer(), closer, nil
}
