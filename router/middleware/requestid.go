package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/1024casts/snake/pkg/utils"
)

// RequestID 透传Request-ID，如果没有则生成一个
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestID := c.Request.Header.Get(utils.XRequestID)

		// Create request id with UUID
		if requestID == "" {
			requestID = utils.GenRequestID()
		}

		// Expose it for use in the application
		c.Set(utils.XRequestID, requestID)

		// Set X-Request-ID header
		c.Writer.Header().Set(utils.XRequestID, requestID)
		c.Next()
	}
}
