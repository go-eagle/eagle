package middleware

import (
	"github.com/1024casts/snake/pkg/net/tracing"
	"github.com/1024casts/snake/pkg/snake"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

var ParentSpan opentracing.Span

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer, closer := tracing.Init(snake.App.Conf.App.Name)
		defer closer.Close()

		spCtx, err := tracing.Extract(tracer, c.Request)
		if err != nil {
			ParentSpan = tracer.StartSpan(c.Request.URL.Path)
			defer ParentSpan.Finish()
		} else {
			ParentSpan = tracer.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(spCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				ext.SpanKindRPCServer,
			)
			defer ParentSpan.Finish()
		}
		c.Next()
	}
}
