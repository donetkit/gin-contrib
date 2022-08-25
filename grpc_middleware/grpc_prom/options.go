package grpc_prom

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type RequestLabelMappingFn func(c *gin.Context) string

// Config defines the config for logger middleware
type config struct {
	handlerUrl              string
	namespace               string
	name                    string
	duration                []float64
	slowTime                float64
	excludeRegexCode        []string
	excludeRegexRpcType     []string
	excludeRegexServiceName []string
	excludeRegexMethodName  []string
	endpointLabelMappingFn  RequestLabelMappingFn

	counterOpts []CounterOption
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

// WithExcludeRegexRpcType set excludeRegexRpcType function regexp
func WithExcludeRegexRpcType(excludeRegexRpcType []string) Option {
	return func(cfg *config) {
		cfg.excludeRegexRpcType = excludeRegexRpcType
	}
}

// WithExcludeRegexRegexServiceName set excludeRegexServiceName function regexp
func WithExcludeRegexRegexServiceName(excludeRegexServiceName []string) Option {
	return func(cfg *config) {
		cfg.excludeRegexServiceName = excludeRegexServiceName
	}
}

// WithExcludeRegexMethodName set excludeRegexMethodName function regexp
func WithExcludeRegexMethodName(excludeRegexMethodName []string) Option {
	return func(cfg *config) {
		cfg.excludeRegexMethodName = excludeRegexMethodName
	}
}

// WithExcludeRegexCode set excludeRegexCode function regexp
func WithExcludeRegexCode(excludeRegexCode []string) Option {
	return func(cfg *config) {
		cfg.excludeRegexCode = excludeRegexCode
	}
}

// WithEndpointLabelMappingFn set endpointLabelMappingFn function
func WithEndpointLabelMappingFn(endpointLabelMappingFn RequestLabelMappingFn) Option {
	return func(cfg *config) {
		cfg.endpointLabelMappingFn = endpointLabelMappingFn
	}
}

//WithPromHandler set router function
func WithPromHandler(r *http.ServeMux) Option {
	return func(cfg *config) {
		if r != nil {
			r.Handle(cfg.handlerUrl, promhttp.Handler())
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

// WithCounterOption set counterOpts function 1
func WithCounterOption(counterOpts []CounterOption) Option {
	return func(cfg *config) {
		cfg.counterOpts = counterOpts
	}
}
