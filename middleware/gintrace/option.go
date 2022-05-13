package gintrace

import (
	"github.com/donetkit/gin-contrib/trace"
)

type config struct {
	tracerServer           *trace.Server
	excludeRegexStatus     []string
	excludeRegexEndpoint   []string
	excludeRegexMethod     []string
	endpointLabelMappingFn RequestLabelMappingFn
}

// Option specifies instrumentation configuration options.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

// WithTracer  tracerServer trace.Server
func WithTracer(tracerServer *trace.Server) Option {
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
