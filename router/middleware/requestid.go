package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/1024casts/snake/pkg/utils"
)

// RequestID 透传Request-ID，如果没有则生成一个
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestID := c.Request.Header.Get(utils.HTTPHeaderXRequestIDKey)

		// Create request id with UUID
		if requestID == "" {
			requestID = utils.GenRequestID()
			ctx := utils.NewRequestIDContext(c.Request.Context(), requestID)
			c.Request = c.Request.WithContext(ctx)
		}

		// Expose it for use in the application
		c.Set(utils.ContextRequestIDKey, requestID)

		// Set X-Request-ID header
		c.Writer.Header().Set(utils.HTTPHeaderXRequestIDKey, requestID)

		c.Next()
	}
}
