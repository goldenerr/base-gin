package log

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"net/http"
)

const TraceContextKey string = "traceContext"

// 生成全局tracer，链接agent，如果agent不可用，使用logger report
func InitTracer() {
	var report jaeger.Reporter
	report = jaeger.NewNullReporter()
	//transport, err := jaeger.NewUDPTransport("jaeger-agent:5775", 0)
	//if err != nil {
	//	report = jaeger.NewNullReporter()
	//} else {
	//	report = jaeger.NewRemoteReporter(transport)
	//}

	tracer, _ := jaeger.NewTracer("todo GetServiceName", jaeger.NewConstSampler(true), report, jaeger.TracerOptions.CustomHeaderKeys(&jaeger.HeadersConfig{
		JaegerDebugHeader:        jaeger.JaegerDebugHeader,
		JaegerBaggageHeader:      jaeger.JaegerBaggageHeader,
		TraceContextHeaderName:   "trace-id",
		TraceBaggageHeaderPrefix: "jhy-",
	}))

	opentracing.SetGlobalTracer(tracer)
	return
}

func TracerStartSpan(parent opentracing.SpanContext, operationName string, tags map[string]interface{}) opentracing.Span {
	var options []opentracing.StartSpanOption
	if parent != nil {
		options = append(options, opentracing.ChildOf(parent))
	}

	span := opentracing.StartSpan(operationName, options...)
	for k, v := range tags {
		span.SetTag(k, v)
	}
	return span
}

// use trace id as request id
func GetRequestIdFromTrace(ctx *gin.Context) string {
	var spanContext jaeger.SpanContext
	if ctx != nil {
		tempSpan, exist := ctx.Get(TraceContextKey)
		if exist {
			spanContext = tempSpan.(opentracing.Span).Context().(jaeger.SpanContext)
			return spanContext.String()
		}
	}
	return ""
}

type traceLogger struct{}

func (logger *traceLogger) Error(msg string) {
	ErrorLogger(nil, msg)
}

func (logger *traceLogger) Infof(msg string, args ...interface{}) {
	DebugfLogger(nil, msg, args...)
}

func Tracer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span := TracerStartSpan(getSpanContextFromRequest(ctx.Request), "HTTP_"+ctx.Request.URL.Path, getSpanTagsByHttpServer(ctx.Request))
		defer span.Finish()
		ctx.Set(TraceContextKey, span)
		ctx.Set("requestId", GetRequestIdFromTrace(ctx))

		ctx.Next()

		render, exist := ctx.Get("render")
		retCode := 0
		if exist {
			retCode = render.(TRenderJson).ReturnCode
		}
		span.SetTag("retCode", retCode)
		span.SetTag(string(ext.HTTPStatusCode), ctx.Writer.Status())
	}
}

func TracerBeforeRun(ctx *gin.Context) {
	ctx.Set("requestId", GetRequestIdFromTrace(ctx))
}

func TracerAfterRun(ctx *gin.Context) {
	retCode := 0
	//if ctx.CustomContext.Error != nil {
	//	retCode = -1
	//}

	tempSpan, exist := ctx.Get(TraceContextKey)
	if exist {
		span := tempSpan.(opentracing.Span)
		span.SetTag("retCode", retCode)
		span.Finish()
	}
}

func getSpanContextFromRequest(req *http.Request) (spanContext opentracing.SpanContext) {
	if req.Header != nil {
		spanContext, _ = opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	}
	return spanContext
}

func getSpanTagsByHttpServer(req *http.Request) map[string]interface{} {
	return map[string]interface{}{
		ext.SpanKindRPCServer.Key: ext.SpanKindRPCServer.Value,
		string(ext.HTTPMethod):    req.Method,
		string(ext.HTTPUrl):       req.URL.Path,
	}
}
