package utils

import (
	"context"

	"github.com/google/uuid"
)

const (
	// ContextRequestIDKey context request id for context
	ContextRequestIDKey = "request_id"

	// HTTPHeaderXRequestIDKey http header request ID key
	HTTPHeaderXRequestIDKey = "X-Request-ID"
)

// GenRequestID 生成随机字符串，eg: 76d27e8c-a80e-48c8-ad20-e5562e0f67e4
func GenRequestID() string {
	reqID, _ := uuid.NewRandom()
	return reqID.String()
}

// GetRequestIDContext 注入到全局context
func GetRequestIDContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextRequestIDKey, requestID)
}

// GetRequestID will get request id from a http request and return it as a string
func GetRequestID(ctx context.Context) string {
	reqID := ctx.Value(ContextRequestIDKey)
	if requestID, ok := reqID.(string); ok {
		return requestID
	}
	return ""
}
