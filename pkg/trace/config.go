package trace

import (
	"fmt"

	"github.com/1024casts/snake/pkg/trace/elastic"
	"github.com/1024casts/snake/pkg/trace/jaeger"
	"github.com/1024casts/snake/pkg/trace/zipkin"
)

var (
	supportedTraceAgent = map[string]bool{
		zipkin.Name:  true,
		jaeger.Name:  true,
		elastic.Name: true,
	}
)

type Config struct {
	ServiceName string // The name of this service
	TraceAgent  string // The type of trace agent: zipkin, jaeger or elastic
	OpenDebug   bool

	Zipkin  zipkin.Config  // Settings for zipkin, only useful when TraceAgent is zipkin
	Jaeger  jaeger.Config  // Settings for jaeger, only useful when TraceAgent is jaeger
	Elastic elastic.Config // Settings for elastic, only useful when TraceAgent is elastic
}

func (cfg *Config) Check() error {

	if len(cfg.TraceAgent) == 0 {
		return fmt.Errorf("ModTrace.TraceAgent not set")
	}

	if _, ok := supportedTraceAgent[cfg.TraceAgent]; !ok {
		return fmt.Errorf("Trace.TraceAgent %s is not supported", cfg.TraceAgent)
	}

	return nil
}

func (cfg *Config) GetTraceConfig() TraceAgent {
	switch cfg.TraceAgent {
	case jaeger.Name:
		return &cfg.Jaeger
	case zipkin.Name:
		return &cfg.Zipkin
	case elastic.Name:
		return &cfg.Elastic
	default:
		return &cfg.Jaeger
	}
}
