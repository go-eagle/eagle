package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

var _ HistogramVec = (*promHistogramVec)(nil)

// HistogramVecOpts is histogram vector opts.
type HistogramVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
	Buckets   []float64
}

// Histogram prom histogram collection.
type promHistogramVec struct {
	histogram *prometheus.HistogramVec
}

// NewHistogramVec new a histogram vec.
func NewHistogramVec(cfg *HistogramVecOpts) HistogramVec {
	if cfg == nil {
		return nil
	}
	vec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.Name,
			Help:      cfg.Help,
			Buckets:   cfg.Buckets,
		}, cfg.Labels)
	prometheus.MustRegister(vec)
	return &promHistogramVec{
		histogram: vec,
	}
}

// Timing adds a single observation to the histogram.
func (histogram *promHistogramVec) Observe(v int64, labels ...string) {
	histogram.histogram.WithLabelValues(labels...).Observe(float64(v))
}
