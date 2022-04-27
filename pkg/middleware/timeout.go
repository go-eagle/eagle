package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	timeout "github.com/vearne/gin-timeout"
)

// Timeout 超时中间件
func Timeout(t time.Duration) gin.HandlerFunc {
	// see:
	// https://github.com/vearne/gin-timeout
	// https://vearne.cc/archives/39135
	// https://github.com/gin-contrib/timeout
	return timeout.Timeout(
		timeout.WithTimeout(t),
	)
}
