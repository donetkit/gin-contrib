package gintrace

import (
	"fmt"
	"github.com/donetkit/contrib/tracer"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"regexp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	tracerKey = "go-contrib-tracer-key"
)

type RequestLabelMappingFn func(c *gin.Context) string

// New returns middleware that will tracer incoming requests.
// The service parameter should describe the name of the (virtual)
// webserve handling the request.
func New(opts ...Option) gin.HandlerFunc {
	cfg := config{
		tracerName: "Service",
		traceIdKey: "trace-id",
		spanIdKey:  "span-id",
		endpointLabelMappingFn: func(c *gin.Context) string {
			return c.Request.URL.Path
		}}
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	return func(c *gin.Context) {
		if cfg.tracerServer == nil {
			return
		}
		endpoint := cfg.endpointLabelMappingFn(c)
		method := c.Request.Method
		isOk := cfg.checkLabel(fmt.Sprintf("%d", c.Writer.Status()), cfg.excludeRegexStatus) && cfg.checkLabel(endpoint, cfg.excludeRegexEndpoint) && cfg.checkLabel(method, cfg.excludeRegexMethod)
		if !isOk {
			return
		}
		c.Set(tracerKey, cfg.tracerServer)
		savedCtx := c.Request.Context()
		defer func() {
			c.Request = c.Request.WithContext(savedCtx)
		}()
		ctx := cfg.tracerServer.Propagators.Extract(savedCtx, propagation.HeaderCarrier(c.Request.Header))
		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", c.Request)...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(c.Request)...),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(cfg.tracerName, c.FullPath(), c.Request)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}

		if values, ok := c.Request.Header["Referer"]; ok && len(values) > 0 {
			opts = append(opts, oteltrace.WithAttributes(attribute.String("http.referer", values[0])))
		}

		if values, ok := c.Request.Header["X-Forwarded-For"]; ok && len(values) > 0 {
			opts = append(opts, oteltrace.WithAttributes(attribute.String("X-Forwarded-For", values[0])))
		}
		
		spanName := c.FullPath()
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", c.Request.Method)
		}
		ctx, span := cfg.tracerServer.Tracer.Start(ctx, spanName, opts...)
		if !span.IsRecording() {
			return
		}
		defer span.End()
		// header写入trace-id和span-id
		if cfg.writerTraceId {
			c.Header(cfg.traceIdKey, span.SpanContext().TraceID().String())
		}
		if cfg.writerSpanId {
			c.Header(cfg.spanIdKey, span.SpanContext().SpanID().String())
		}
		// pass the span through the request context
		c.Request = c.Request.WithContext(ctx)
		// serve the request to the next middleware
		c.Next()
		status := c.Writer.Status()
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)
		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)
		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.String("gin.errors", c.Errors.String()))
		}
	}
}

// HTML will tracer the rendering of the template as a child of the
// span in the given context. This is a replacement for
// gin.Context.HTML function - it invokes the original function after
// setting up the span.
func HTML(c *gin.Context, code int, name string, obj interface{}) {
	var trace tracer.Server
	tracerInterface, ok := c.Get(tracerKey)
	if ok {
		trace, ok = tracerInterface.(tracer.Server)
	}
	if !ok {
		return
	}
	savedContext := c.Request.Context()
	defer func() {
		c.Request = c.Request.WithContext(savedContext)
	}()
	opt := oteltrace.WithAttributes(attribute.String("go.template", name))
	_, span := trace.Tracer.Start(savedContext, "gin.renderer.html", opt)
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("error rendering template:%s: %s", name, r)
			span.RecordError(err)
			span.SetStatus(codes.Error, "template failure")
			span.End()
			panic(r)
		} else {
			span.End()
		}
	}()
	c.HTML(code, name, obj)
}

// checkLabel returns the match result of labels.
// Return true if regex-pattern compiles failed.
func (c *config) checkLabel(label string, patterns []string) bool {
	if len(patterns) <= 0 {
		return true
	}
	for _, pattern := range patterns {
		if pattern == "" {
			return true
		}
		matched, err := regexp.MatchString(pattern, label)
		if err != nil {
			return true
		}
		if matched {
			return false
		}
	}
	return true
}
