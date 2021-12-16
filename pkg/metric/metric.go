package metric

// CounterVec counter vec.
type CounterVec interface {
	// Inc increments the counter by 1. Use Add to increment it by arbitrary
	// non-negative values.
	Inc(labels ...string)
	// Add adds the given value to the counter. It panics if the value is <
	// 0.
	Add(v float64, labels ...string)
}

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

// HistogramVec gauge vec.
type HistogramVec interface {
	// Observe adds a single observation to the histogram.
	Observe(v int64, labels ...string)
}

// VectorOpts contains the common arguments for creating vec Metric..
type VectorOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

// Metric is a sample interface.
// Implementations of Metrics in metric package are Counter, Gauge,
// PointGauge, RollingCounter and RollingGauge.
type Metric interface {
	// Add adds the given value to the counter.
	Add(int64)
	// Value gets the current value.
	// If the metric's type is PointGauge, RollingCounter, RollingGauge,
	// it returns the sum value within the window.
	Value() int64
}

// Aggregation contains some common aggregation function.
// Each aggregation can compute summary statistics of window.
type Aggregation interface {
	// Min finds the min value within the window.
	Min() float64
	// Max finds the max value within the window.
	Max() float64
	// Avg computes average value within the window.
	Avg() float64
	// Sum computes sum value within the window.
	Sum() float64
}
