package requestid

import (
	"context"

	"github.com/google/uuid"
)

const (
	// ContextRequestIDKey context request id for context
	ContextRequestIDKey = "request_id"

	// HeaderXRequestIDKey http header request ID key
	HeaderXRequestIDKey = "X-Request-ID"
)

// Generate 生成随机字符串，eg: 76d27e8c-a80e-48c8-ad20-e5562e0f67e4
func Generate() string {
	reqID, _ := uuid.NewRandom()
	return reqID.String()
}

// NewContext creates a context with request id
func NewContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextRequestIDKey, requestID)
}

// GetFromContext will get request id from a http request and return it as a string
func GetFromContext(ctx context.Context) string {
	reqID := ctx.Value(ContextRequestIDKey)
	if requestID, ok := reqID.(string); ok {
		return requestID
	}
	return ""
}
