package trace

// Config jaeger config
type Config struct {
	ServiceName        string // The name of this service
	LocalAgentHostPort string
	CollectorEndpoint  string
}
