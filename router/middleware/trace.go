package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/1024casts/snake/pkg/net/tracing"
)

const (
	// DefaultServiceName service name
	DefaultServiceName = "snake"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer, closer := tracing.Init(DefaultServiceName)
		defer closer.Close()

		// set into opentracing
		opentracing.SetGlobalTracer(tracer)

		var sp opentracing.Span
		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		ctx, _ := tracer.Extract(opentracing.HTTPHeaders, carrier)
		sp = tracer.StartSpan(c.Request.URL.Path, ext.RPCServerOption(ctx))
		// or use under statement
		// sp = tracing.StartSpanFromRequest(c.Request.URL.Path, tracer, c.Request)
		defer sp.Finish()

		// record HTTP method
		ext.HTTPMethod.Set(sp, c.Request.Method)
		// record HTTP url
		ext.HTTPUrl.Set(sp, c.Request.URL.String())
		// record component name
		// ext.Component.Set(sp, componentName)

		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), sp))

		c.Next()

		// record HTTP status code
		ext.HTTPStatusCode.Set(sp, uint16(c.Writer.Status()))
	}
}
