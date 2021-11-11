package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-eagle/eagle/pkg/app"

	"github.com/gin-gonic/gin"
)

// Timeout 超时中间件
func Timeout(t time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		defer cancel()

		defer func() {
			// check if context timeout was reached
			if ctx.Err() == context.DeadlineExceeded {
				// write response and abort the request
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.AbortWithStatusJSON(http.StatusGatewayTimeout, app.Response{
					Code:    http.StatusGatewayTimeout,
					Message: ctx.Err().Error(),
				})
			}
		}()

		// replace request with context wrapped request
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
