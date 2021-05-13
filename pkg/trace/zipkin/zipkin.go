package zipkin

import (
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter/http"
)

// Name sets the name of this tracer.
const Name = "zipkin"

// Config provides configuration settings for a zipkin tracer.
type Config struct {
	HTTPEndpoint string  // HTTP Endpoint to report traces to
	SameSpan     bool    // Use Zipkin SameSpan RPC style traces
	ID128Bit     bool    // Use Zipkin 128 bit root span IDs
	SampleRate   float64 // The rate between 0.0 and 1.0 of requests to trace
}

// SetDefaults sets the default values.
func (c *Config) SetDefaults() {
	c.HTTPEndpoint = "http://localhost:9411/api/v2/spans"
	c.SameSpan = false
	c.ID128Bit = true
	c.SampleRate = 1.0
}

// New sets up the tracer
func (c *Config) New(serviceName string) (opentracing.Tracer, io.Closer, error) {
	// create our local endpoint
	endpoint, err := zipkin.NewEndpoint(serviceName, "0.0.0.0:0")
	if err != nil {
		return nil, nil, err
	}

	// create our sampler
	sampler, err := zipkin.NewBoundarySampler(c.SampleRate, time.Now().Unix())
	if err != nil {
		return nil, nil, err
	}

	// create the span reporter
	reporter := http.NewReporter(c.HTTPEndpoint)

	// create the native Zipkin tracer
	nativeTracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithSharedSpans(c.SameSpan),
		zipkin.WithTraceID128Bit(c.ID128Bit),
		zipkin.WithSampler(sampler),
	)
	if err != nil {
		return nil, nil, err
	}

	// wrap the Zipkin native tracer with the OpenTracing Bridge
	tracer := zipkinot.Wrap(nativeTracer)

	// Without this, child spans are getting the NOOP tracer
	opentracing.SetGlobalTracer(tracer)

	return tracer, reporter, nil
}
