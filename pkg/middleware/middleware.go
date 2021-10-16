package middleware

import "github.com/gin-gonic/gin"

// Middlewares global middleware
var Middlewares = defaultMiddlewares()

func defaultMiddlewares() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"recovery":   gin.Recovery(),
		"secure":     Secure,
		"options":    Options,
		"nocache":    NoCache,
		"logger":     Logging(),
		"cors":       Cors(),
		"request_id": RequestID(),
	}
}
