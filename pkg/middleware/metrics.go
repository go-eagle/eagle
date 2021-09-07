package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/go-eagle/eagle/pkg/conf"

	"github.com/go-eagle/eagle/pkg/metric"

	"github.com/gin-gonic/gin"
)

var namespace = conf.Conf.App.Name

var (
	labels = []string{"status", "endpoint", "method"}

	// 应用时长
	uptime = metric.NewCounterVec(
		&metric.CounterVecOpts{
			Namespace: namespace,
			Name:      "uptime",
			Help:      "HTTP service uptime.",
			Labels:    labels,
		})

	// QPS
	reqCount = metric.NewCounterVec(
		&metric.CounterVecOpts{
			Namespace: namespace,
			Name:      "http_request_count_total",
			Help:      "Total number of HTTP requests made.",
			Labels:    labels,
		})

	// 接口响应时间
	reqDuration = metric.NewHistogramVec(
		&metric.HistogramVecOpts{
			Namespace: namespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request latencies in seconds.",
			Labels:    labels,
		})

	// 请求大小
	reqSizeBytes = metric.NewHistogramVec(
		&metric.HistogramVecOpts{
			Namespace: namespace,
			Name:      "http_request_size_bytes",
			Help:      "HTTP request sizes in bytes.",
			Labels:    labels,
		})

	// 响应大小
	respSizeBytes = metric.NewHistogramVec(
		&metric.HistogramVecOpts{
			Namespace: namespace,
			Name:      "http_response_size_bytes",
			Help:      "HTTP request sizes in bytes.",
			Labels:    labels,
		})
)

// init registers the prometheus metrics
func init() {
	// go recordUptime()
}

// recordUptime increases service uptime per second.
func recordUptime() {
	for range time.Tick(time.Second) {
		uptime.Inc()
	}
}

// calcRequestSize returns the size of request object.
func calcRequestSize(r *http.Request) float64 {
	size := 0
	if r.URL != nil {
		size = len(r.URL.String())
	}

	size += len(r.Method)
	size += len(r.Proto)

	for name, values := range r.Header {
		size += len(name)
		for _, value := range values {
			size += len(value)
		}
	}
	size += len(r.Host)

	// r.Form and r.MultipartForm are assumed to be included in r.URL.
	if r.ContentLength != -1 {
		size += int(r.ContentLength)
	}
	return float64(size)
}

type RequestLabelMappingFn func(c *gin.Context) string

// PromOpts represents the Prometheus middleware Options.
// It is used for filtering labels by regex.
type PromOpts struct {
	ExcludeRegexStatus     string
	ExcludeRegexEndpoint   string
	ExcludeRegexMethod     string
	EndpointLabelMappingFn RequestLabelMappingFn
}

// NewDefaultOpts return the default ProOpts
func NewDefaultOpts() *PromOpts {
	return &PromOpts{
		EndpointLabelMappingFn: func(c *gin.Context) string {
			//by default do nothing, return URL as is
			return c.Request.URL.Path
		},
	}
}

// checkLabel returns the match result of labels.
// Return true if regex-pattern compiles failed.
func (po *PromOpts) checkLabel(label, pattern string) bool {
	if pattern == "" {
		return true
	}

	matched, err := regexp.MatchString(pattern, label)
	if err != nil {
		return true
	}
	return !matched
}

// Metrics returns a gin.HandlerFunc for exporting some Web metrics
func Metrics(promOpts *PromOpts) gin.HandlerFunc {
	// make sure promOpts is not nil
	if promOpts == nil {
		promOpts = NewDefaultOpts()
	}

	// make sure EndpointLabelMappingFn is callable
	if promOpts.EndpointLabelMappingFn == nil {
		promOpts.EndpointLabelMappingFn = func(c *gin.Context) string {
			return c.Request.URL.Path
		}
	}

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := fmt.Sprintf("%d", c.Writer.Status())
		endpoint := promOpts.EndpointLabelMappingFn(c)
		method := c.Request.Method

		labels := []string{status, endpoint, method}

		isOk := promOpts.checkLabel(status, promOpts.ExcludeRegexStatus) &&
			promOpts.checkLabel(endpoint, promOpts.ExcludeRegexEndpoint) &&
			promOpts.checkLabel(method, promOpts.ExcludeRegexMethod)

		if !isOk {
			return
		}
		// no response content will return -1
		respSize := c.Writer.Size()
		if respSize < 0 {
			respSize = 0
		}
		reqCount.Inc(labels...)
		reqDuration.Observe(int64(time.Since(start).Seconds()), labels...)
		reqSizeBytes.Observe(int64(calcRequestSize(c.Request)), labels...)
		respSizeBytes.Observe(int64(respSize), labels...)
	}
}

// PromHandler wrappers the standard http.Handler to gin.HandlerFunc
func PromHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
