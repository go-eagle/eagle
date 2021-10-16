package trace

import (
	"fmt"

	"github.com/go-eagle/eagle/pkg/trace/jaeger"
)

var (
	supportedTraceAgent = map[string]bool{
		jaeger.Name: true,
	}
)

// Config .
type Config struct {
	ServiceName string // The name of this service
	TraceAgent  string // The type of trace agent: zipkin, jaeger or elastic
	OpenDebug   bool

	Jaeger jaeger.Config // Settings for jaeger, only useful when TraceAgent is jaeger
}

// Check check config
func (cfg *Config) Check() error {
	if len(cfg.TraceAgent) == 0 {
		return fmt.Errorf("ModTrace.TraceAgent not set")
	}

	if _, ok := supportedTraceAgent[cfg.TraceAgent]; !ok {
		return fmt.Errorf("Trace.TraceAgent %s is not supported", cfg.TraceAgent)
	}

	return nil
}
