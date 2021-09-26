package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// ContextRequestIDKey context request id for context
	ContextRequestIDKey = "request_id"

	// HeaderXRequestIDKey http header request ID key
	HeaderXRequestIDKey = "X-Request-ID"
)

// RequestID is a middleware that injects a 'X-Request-ID' into the context and request/response header of each request.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestID := c.Request.Header.Get(HeaderXRequestIDKey)

		// Create request id with UUID
		if requestID == "" {
			requestID = generateID()
			c.Request.Header.Set(HeaderXRequestIDKey, requestID)
			// Expose it for use in the application
			c.Set(ContextRequestIDKey, requestID)
		}

		// Set X-Request-ID header
		c.Writer.Header().Set(HeaderXRequestIDKey, requestID)

		c.Next()
	}
}

// generateID 生成随机字符串，eg: 76d27e8c-a80e-48c8-ad20-e5562e0f67e4
func generateID() string {
	reqID, _ := uuid.NewRandom()
	return reqID.String()
}

// GetRequestIDFromContext returns 'RequestID' from the given context if present.
func GetRequestIDFromContext(c *gin.Context) string {
	if v, ok := c.Get(HeaderXRequestIDKey); ok {
		if requestID, ok := v.(string); ok {
			return requestID
		}
	}

	return ""
}

// GetRequestIDFromHeaders returns 'RequestID' from the headers if present.
func GetRequestIDFromHeaders(c *gin.Context) string {
	return c.Request.Header.Get(HeaderXRequestIDKey)
}
