package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// AccessLog record access log.
func AccessLog() gin.HandlerFunc {
	return gin.LoggerWithConfig(
		gin.LoggerConfig{
			Formatter: func(params gin.LogFormatterParams) string {
				return fmt.Sprintf("%s - [%s] \"%s %s %s\" %d %d %s \"%s\" \"%s\"\n",
					params.ClientIP,
					params.TimeStamp.Format(time.RFC1123Z),
					params.Method,
					params.Path,
					params.Request.Proto,
					params.StatusCode,
					params.BodySize,
					params.Latency,
					params.Request.UserAgent(),
					params.ErrorMessage,
				)
			},
			SkipPaths: []string{"/health", "/metrics"},
		})
}
