package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/1024casts/snake/pkg/net/tracing"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer, closer := tracing.Init("snake")
		defer closer.Close()

		// set into opentracing
		opentracing.SetGlobalTracer(tracer)

		parentSpanCtx, err := tracing.Extract(tracer, c.Request)
		if err != nil {
			parentSpan := tracer.StartSpan(c.Request.URL.Path)
			defer parentSpan.Finish()
		} else {
			childSpan := tracer.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(parentSpanCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				ext.SpanKindRPCServer,
			)
			defer childSpan.Finish()
		}
		c.Next()
	}
}
