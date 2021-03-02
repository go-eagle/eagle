package redis

import (
	"context"
	"os"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

var redisAddr string
var client *redis.Client
var tracer *mocktracer.MockTracer

func init() {
	tracer = mocktracer.New()
	opentracing.SetGlobalTracer(tracer)
}

func TestMain(m *testing.M) {
	// in-memory redis
	miniRedis, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer miniRedis.Close()

	redisAddr = miniRedis.Addr()

	client = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	os.Exit(m.Run())
}

// SET

func TestSet(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	span, ctx := opentracing.StartSpanFromContext(ctx, "test-set")
	ctxClient := WrapRedisClient(ctx, client)
	callSet(assert, ctxClient, "with span")
	span.Finish()

	spans := tracer.FinishedSpans()
	assert.Len(spans, 2, "the number of finished spans is invalid")

	redisSpan := spans[0]
	assert.Equal("redis-cmd", redisSpan.OperationName)

	expectedTags := buildExpectedTags("set")
	assertTags(assert, redisSpan, expectedTags)

	tracer.Reset()
}

func TestSetPipeline(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	span, ctx := opentracing.StartSpanFromContext(ctx, "test-set-pipeline")
	ctxClient := WrapRedisClient(ctx, client)
	setPipelineParams := make(map[string]string)
	setPipelineParams["foo"] = "with span on foo pipeline"
	setPipelineParams["bar"] = "with span on bar pipeline"
	callSetPipeline(assert, ctxClient, setPipelineParams)
	span.Finish()

	spans := tracer.FinishedSpans()
	assert.Len(spans, 2, "the number of finished spans is invalid")

	redisSpan := spans[0]
	assert.Equal("redis-pipeline-cmd", redisSpan.OperationName)

	expectedTags := buildExpectedTags("set -> set")
	assertTags(assert, redisSpan, expectedTags)

	tracer.Reset()
}

func callSet(assert *assert.Assertions, client *redis.Client, value string) {
	_, err := client.Set("foo", value, 0).Result()
	assert.Nil(err, "Redis returned error: %v", err)
}

func callSetPipeline(assert *assert.Assertions, client *redis.Client, setPipelineParams map[string]string) {
	pipeline := client.Pipeline()
	for key, value := range setPipelineParams {
		pipeline.Set(key, value, 0)
	}
	_, err := pipeline.Exec()
	assert.Nil(err, "Redis returned error: %v", err)
}

// GET

func TestGet(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	span, ctx := opentracing.StartSpanFromContext(ctx, "test-get")
	ctxClient := WrapRedisClient(ctx, client)
	callGet(assert, ctxClient)
	span.Finish()

	spans := tracer.FinishedSpans()
	assert.Len(spans, 2, "the number of finished spans is invalid")

	redisSpan := spans[0]
	assert.Equal("redis-cmd", redisSpan.OperationName)

	expectedTags := buildExpectedTags("get")
	assertTags(assert, redisSpan, expectedTags)

	tracer.Reset()
}

func TestGetPipeline(t *testing.T) {
	ctx := context.Background()
	assert := assert.New(t)

	span, ctx := opentracing.StartSpanFromContext(ctx, "test-get-pipeline")
	ctxClient := WrapRedisClient(ctx, client)
	getPipelineParams := []string{"foo", "bar"}
	callGetPipeline(assert, ctxClient, getPipelineParams)
	span.Finish()

	spans := tracer.FinishedSpans()
	assert.Len(spans, 2, "the number of finished spans is invalid")

	redisSpan := spans[0]
	assert.Equal("redis-pipeline-cmd", redisSpan.OperationName)

	expectedTags := buildExpectedTags("get -> get")
	assertTags(assert, redisSpan, expectedTags)

	tracer.Reset()
}

func callGet(assert *assert.Assertions, client *redis.Client) {
	_, err := client.Get("foo").Result()
	assert.Nil(err, "Redis returned error: %v", err)
}

func callGetPipeline(assert *assert.Assertions, client *redis.Client, getPipelineParams []string) {
	pipeline := client.Pipeline()
	for _, key := range getPipelineParams {
		pipeline.Get(key)
	}
	_, err := pipeline.Exec()
	assert.Nil(err, "Redis returned error: %v", err)
}

// MISC

func buildExpectedTags(expectedDbMethod string) map[string]interface{} {
	expectedTags := make(map[string]interface{})
	expectedTags["db.type"] = "redis"
	expectedTags["db.method"] = expectedDbMethod
	expectedTags["peer.address"] = redisAddr
	expectedTags["span.kind"] = ext.SpanKindEnum("client")
	return expectedTags
}

func assertTags(assert *assert.Assertions, redisSpan *mocktracer.MockSpan, expectedTags map[string]interface{}) map[string]interface{} {
	actualTags := redisSpan.Tags()
	assert.Len(actualTags, len(expectedTags), "redis span tags number is invalid")
	for expectedTagKey, expectedTagValue := range expectedTags {
		actualTagValue, ok := actualTags[expectedTagKey]
		assert.True(ok, "redis span doesn't have tag '%s'", expectedTagKey)
		assert.Equal(expectedTagValue, actualTagValue, "redis span tag '%s' is invalid", expectedTagKey)
	}
	return actualTags
}
