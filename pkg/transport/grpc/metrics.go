package grpc

import "github.com/prometheus/client_golang/prometheus"

var (
	namespace = "grpc"

	sentBytes = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "network",
		Name:      "client_grpc_sent_bytes_total",
		Help:      "The total number of bytes sent to grpc clients.",
	})

	receivedBytes = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "network",
		Name:      "client_grpc_received_bytes_total",
		Help:      "The total number of bytes received from grpc clients.",
	})

	streamFailures = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "network",
		Name:      "server_stream_failures_total",
		Help:      "The total number of stream failures from the local server.",
	},
		[]string{"Type", "API"},
	)

	clientRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "client_requests_total",
		Help:      "The total number of client requests per client version.",
	},
		[]string{"type", "client_api_version"},
	)
)

func init() {
	prometheus.MustRegister(sentBytes)
	prometheus.MustRegister(receivedBytes)
	prometheus.MustRegister(streamFailures)
	prometheus.MustRegister(clientRequests)
}
