package gintrace

import (
	"github.com/donetkit/contrib/tracer"
)

type config struct {
	tracerName             string
	tracerServer           *tracer.Server
	excludeRegexStatus     []string
	excludeRegexEndpoint   []string
	excludeRegexMethod     []string
	endpointLabelMappingFn RequestLabelMappingFn
	writerTraceId          bool
	writerSpanId           bool
	traceIdKey             string
	spanIdKey              string
}

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

// WithTracer  tracerServer tracer.Server
func WithTracer(tracerServer *tracer.Server) Option {
	return optionFunc(func(cfg *config) {
		cfg.tracerServer = tracerServer
	})
}

// WithExcludeRegexMethod set excludeRegexMethod function regexp
func WithExcludeRegexMethod(excludeRegexMethod []string) Option {
	return optionFunc(func(cfg *config) {
		cfg.excludeRegexMethod = excludeRegexMethod
	})
}

// WithExcludeRegexStatus set excludeRegexStatus function regexp
func WithExcludeRegexStatus(excludeRegexStatus []string) Option {
	return optionFunc(func(cfg *config) {
		cfg.excludeRegexStatus = excludeRegexStatus
	})
}

// WithExcludeRegexEndpoint set excludeRegexEndpoint function regexp
func WithExcludeRegexEndpoint(excludeRegexEndpoint []string) Option {
	return optionFunc(func(cfg *config) {
		cfg.excludeRegexEndpoint = excludeRegexEndpoint
	})
}

// WithEndpointLabelMappingFn set endpointLabelMappingFn function
func WithEndpointLabelMappingFn(endpointLabelMappingFn RequestLabelMappingFn) Option {
	return optionFunc(func(cfg *config) {
		cfg.endpointLabelMappingFn = endpointLabelMappingFn
	})
}

// WithWriterTraceId  writer traceId function
func WithWriterTraceId() Option {
	return optionFunc(func(cfg *config) {
		cfg.writerTraceId = true
	})
}

// WithWriterSpanId  writer spanId function
func WithWriterSpanId() Option {
	return optionFunc(func(cfg *config) {
		cfg.writerSpanId = true
	})
}

// WithName  name default Service
func WithName(name string) Option {
	return optionFunc(func(cfg *config) {
		cfg.tracerName = name
	})
}

// WithTraceIdKey  traceIdKey default trace-id
func WithTraceIdKey(traceIdKey string) Option {
	return optionFunc(func(cfg *config) {
		cfg.traceIdKey = traceIdKey
	})
}

// WithSpanIdKey  spanIdKey default span-id
func WithSpanIdKey(spanIdKey string) Option {
	return optionFunc(func(cfg *config) {
		cfg.spanIdKey = spanIdKey
	})
}
