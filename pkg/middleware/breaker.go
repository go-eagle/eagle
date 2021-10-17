package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/aegis/circuitbreaker"
	"github.com/go-kratos/aegis/circuitbreaker/sre"
	"github.com/pkg/errors"

	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/container/group"
)

// ErrNotAllowed is request failed due to circuit breaker triggered.
var ErrNotAllowed = errors.New("[BREAKER] request failed due to circuit breaker is open")

// BreakerOption is circuit breaker option.
type BreakerOption func(*options)

// WithGroup with circuit breaker group.
// NOTE: implements generics circuitbreaker.CircuitBreaker
func WithGroup(g *group.Group) BreakerOption {
	return func(o *options) {
		o.group = g
	}
}

type options struct {
	group *group.Group
}

// Breaker a circuit breaker middleware
func Breaker(opts ...BreakerOption) gin.HandlerFunc {
	opt := &options{
		group: group.NewGroup(func() interface{} {
			return sre.NewBreaker()
		}),
	}
	for _, o := range opts {
		o(opt)
	}
	return func(c *gin.Context) {
		breaker := opt.group.Get(c.FullPath()).(circuitbreaker.CircuitBreaker)
		if err := breaker.Allow(); err != nil {
			// rejected
			// NOTE: when client reject requets locally,
			// continue add counter let the drop ratio higher.
			breaker.MarkFailed()
			app.NewResponse().Error(c, ErrNotAllowed)
			c.Abort()
			return
		}
		// allowed
		code := c.Writer.Status()
		// NOTE: need to check internal and service unavailable error
		if code == http.StatusInternalServerError || code == http.StatusServiceUnavailable || code == http.StatusGatewayTimeout {
			breaker.MarkFailed()
		} else {
			breaker.MarkSuccess()
		}

		c.Next()
	}
}
