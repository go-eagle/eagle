package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

// GaugeVecOpts is an alias of VectorOpts.
type GaugeVecOpts VectorOpts

// GaugeVec gauge vec.
type GaugeVec interface {
	// Set sets the Gauge to an arbitrary value.
	Set(v float64, labels ...string)
	// Inc increments the Gauge by 1. Use Add to increment it by arbitrary
	// values.
	Inc(labels ...string)
	Dec(labels ...string)
	// Add adds the given value to the Gauge. (The value can be negative,
	// resulting in a decrease of the Gauge.)
	Add(v float64, labels ...string)
}

// gaugeVec gauge vec.
type promGaugeVec struct {
	gauge *prometheus.GaugeVec
}

// NewGaugeVec .
func NewGaugeVec(cfg *GaugeVecOpts) GaugeVec {
	if cfg == nil {
		return nil
	}
	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.Name,
			Help:      cfg.Help,
		}, cfg.Labels)
	prometheus.MustRegister(vec)
	return &promGaugeVec{
		gauge: vec,
	}
}

// Inc increments the counter by 1. Use Inc to increment it by arbitrary.
func (gauge *promGaugeVec) Inc(labels ...string) {
	gauge.gauge.WithLabelValues(labels...).Inc()
}

// Dec increments the counter by 1. Use Dec to increment it by arbitrary.
func (gauge *promGaugeVec) Dec(labels ...string) {
	gauge.gauge.WithLabelValues(labels...).Dec()
}

// Add Inc increments the counter by 1. Use Add to increment it by arbitrary.
func (gauge *promGaugeVec) Add(v float64, labels ...string) {
	gauge.gauge.WithLabelValues(labels...).Add(v)
}

// Set set the given value to the collection.
func (gauge *promGaugeVec) Set(v float64, labels ...string) {
	gauge.gauge.WithLabelValues(labels...).Set(v)
}
