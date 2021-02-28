package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
)

func TestExample(t *testing.T) {

	tracer, closer := jaeger.NewTracer(
		"serviceName",
		jaeger.NewConstSampler(true),
		jaeger.NewInMemoryReporter(),
	)
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	fn := func(c *gin.Context) {
		span := opentracing.SpanFromContext(c.Request.Context())
		if span == nil {
			t.Error("Span is nil")
		}
	}

	r := gin.New()
	r.Use(Trace())
	group := r.Group("")
	group.GET("", fn)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("Error non-nil %v", err)
	}
	r.ServeHTTP(w, req)
}

func TestExampleWithValues(t *testing.T) {

	sampled := "1"
	spanid := "1f369e8cc0105020"
	traceid := "75b353311e9ba7da"
	parentspanid := "2eb00e948dc5066d"

	propagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	tracer, closer := jaeger.NewTracer(
		"serviceName",
		jaeger.NewConstSampler(true),
		jaeger.NewInMemoryReporter(),
		jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, propagator),
		jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, propagator),
		jaeger.TracerOptions.ZipkinSharedRPCSpan(true),
	)
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	fn := func(c *gin.Context) {
		span := opentracing.SpanFromContext(c.Request.Context())
		if span == nil {
			t.Error("Span is nil")
		}
		tracer := opentracing.GlobalTracer()
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Errorf("Error non-nil %v", err)
		}
		carrier := opentracing.HTTPHeadersCarrier(req.Header)
		// Verify the context was populated as expected by the middleware
		err = tracer.Inject(span.Context(), opentracing.HTTPHeaders, carrier)
		if err != nil {
			t.Errorf("Error non-nil %v", err)
		}
		if req.Header.Get("X-B3-Sampled") != sampled {
			t.Errorf("Sampled %s, wanted %s",
				req.Header.Get("X-B3-Sampled"), sampled)
		}
		if req.Header.Get("X-B3-Spanid") != spanid {
			t.Errorf("Spanid %s, wanted %s",
				req.Header.Get("X-B3-Spanid"), spanid)
		}
		if req.Header.Get("X-B3-Traceid") != traceid {
			t.Errorf("Traceid %s, wanted %s",
				req.Header.Get("X-B3-Traceid"), traceid)
		}
		if req.Header.Get("X-B3-Parentspanid") != parentspanid {
			t.Errorf("Parentspanid %s, wanted %s",
				req.Header.Get("X-B3-Parentspanid"), parentspanid)
		}
	}

	r := gin.New()
	r.Use(Trace())
	group := r.Group("")
	group.GET("", fn)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("Error non-nil %v", err)
	}
	req.Header.Set("X-B3-Sampled", sampled)
	req.Header.Set("X-B3-Spanid", spanid)
	req.Header.Set("X-B3-Traceid", traceid)
	req.Header.Set("X-B3-Parentspanid", parentspanid)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}
