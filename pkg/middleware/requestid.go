package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/go-eagle/eagle/pkg/requestid"
)

// RequestID request id middleware
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestID := c.Request.Header.Get(requestid.HeaderXRequestIDKey)

		// Create request id with UUID
		if requestID == "" {
			requestID = requestid.Generate()
			ctx := requestid.NewContext(c.Request.Context(), requestID)
			c.Request = c.Request.WithContext(ctx)
		}

		// Expose it for use in the application
		c.Set(requestid.ContextRequestIDKey, requestID)

		// Set X-Request-ID header
		c.Writer.Header().Set(requestid.HeaderXRequestIDKey, requestID)

		c.Next()
	}
}
