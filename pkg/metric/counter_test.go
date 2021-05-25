package metric

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewCounterVec(t *testing.T) {
	counterVec := NewCounterVec(&CounterVecOpts{
		Namespace: "http_server",
		Subsystem: "requests",
		Name:      "total",
		Help:      "rpc client requests error count.",
	})
	counterVecNil := NewCounterVec(nil)
	assert.NotNil(t, counterVec)
	assert.Nil(t, counterVecNil)
}

func TestCounterIncr(t *testing.T) {
	counterVec := NewCounterVec(&CounterVecOpts{
		Namespace: "http_client",
		Subsystem: "call",
		Name:      "code_total",
		Help:      "http client requests error count.",
		Labels:    []string{"path", "code"},
	})
	cv, _ := counterVec.(*promCounterVec)
	cv.Inc("/v1/users", "500")
	cv.Inc("/v1/users", "500")
	r := testutil.ToFloat64(cv.counter)
	assert.Equal(t, float64(2), r)
}

func TestCounterAdd(t *testing.T) {
	counterVec := NewCounterVec(&CounterVecOpts{
		Namespace: "rpc_server",
		Subsystem: "requests",
		Name:      "err_total",
		Help:      "rpc client requests error count.",
		Labels:    []string{"method", "code"},
	})
	cv, _ := counterVec.(*promCounterVec)
	cv.Add(11, "/v1/users", "500")
	cv.Add(22, "/v1/users", "500")
	r := testutil.ToFloat64(cv.counter)
	assert.Equal(t, float64(33), r)
}
