package logger

import "github.com/donetkit/gin-contrib-log/glog"

// Config defines the config for logger middleware
type config struct {
	// Optional. Default value is gin.defaultLogFormatter
	formatter              LogFormatter
	logger                 glog.ILogger
	excludeRegexStatus     []string
	excludeRegexEndpoint   []string
	excludeRegexMethod     []string
	endpointLabelMappingFn RequestLabelMappingFn
}

// Option for queue system
type Option func(*config)

type Generator func() string

type HeaderStrKey string

// WithLogger set logger function
func WithLogger(logger glog.ILogger) Option {
	return func(cfg *config) {
		cfg.logger = logger
	}
}

// WithExcludeRegexMethod set excludeRegexMethod function regexp
func WithExcludeRegexMethod(excludeRegexMethod []string) Option {
	return func(cfg *config) {
		cfg.excludeRegexMethod = excludeRegexMethod
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

// WithEndpointLabelMappingFn set endpointLabelMappingFn function
func WithEndpointLabelMappingFn(endpointLabelMappingFn RequestLabelMappingFn) Option {
	return func(cfg *config) {
		cfg.endpointLabelMappingFn = endpointLabelMappingFn
	}
}

// WithFormatter set formatter function
func WithFormatter(formatter LogFormatter) Option {
	return func(cfg *config) {
		cfg.formatter = formatter
	}
}
