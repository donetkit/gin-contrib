package prom

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config defines the config for logger middleware
type config struct {
	handlerUrl             string
	namespace              string
	name                   string
	duration               []float64
	slowTime               float64
	excludeRegexStatus     []string
	excludeRegexEndpoint   []string
	excludeRegexMethod     []string
	endpointLabelMappingFn RequestLabelMappingFn
}

// Option for queue system
type Option func(*config)

// WithNamespace set namespace function
func WithNamespace(namespace string) Option {
	return func(cfg *config) {
		cfg.namespace = namespace
	}
}

// WithName set name function
func WithName(name string) Option {
	return func(cfg *config) {
		cfg.name = name
	}
}

// WithHandlerUrl set handlerUrl function
func WithHandlerUrl(handlerUrl string) Option {
	return func(cfg *config) {
		cfg.handlerUrl = handlerUrl
	}
}

// WithExcludeRegexStatus set excludeRegexStatus function regexp
func WithExcludeRegexStatus(excludeRegexStatus []string) Option {
	return func(cfg *config) {
		cfg.excludeRegexStatus = excludeRegexStatus
	}
}

// WithExcludeRegexEndpoint set excludeRegexEndpoint function regexp
func WithExcludeRegexEndpoint(excludeRegexEndpoint []string) Option {
	return func(cfg *config) {
		cfg.excludeRegexEndpoint = excludeRegexEndpoint
	}
}

// WithExcludeRegexMethod set excludeRegexMethod function regexp
func WithExcludeRegexMethod(excludeRegexMethod []string) Option {
	return func(cfg *config) {
		cfg.excludeRegexMethod = excludeRegexMethod
	}
}

// WithEndpointLabelMappingFn set endpointLabelMappingFn function
func WithEndpointLabelMappingFn(endpointLabelMappingFn RequestLabelMappingFn) Option {
	return func(cfg *config) {
		cfg.endpointLabelMappingFn = endpointLabelMappingFn
	}
}

// WithPromHandler set router function
func WithPromHandler(router *gin.Engine) Option {
	return func(cfg *config) {
		if router != nil {
			router.GET(cfg.handlerUrl, promHandler(promhttp.Handler()))
		}
	}
}

// WithDuration set duration function 0.1, 0.3, 1.2, 5
func WithDuration(duration []float64) Option {
	return func(cfg *config) {
		cfg.duration = duration
	}
}

// WithSlowTime set slowTime function 1
func WithSlowTime(slowTime float64) Option {
	return func(cfg *config) {
		cfg.slowTime = slowTime
	}
}
