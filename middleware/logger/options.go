package logger

import "github.com/donetkit/contrib-log/glog"

// Config defines the config for logger middleware
type config struct {
	// Optional. Default value is gin.defaultLogFormatter
	formatter              LogFormatter
	logger                 glog.ILoggerEntry
	excludeRegexStatus     []string
	excludeRegexEndpoint   []string
	excludeRegexMethod     []string
	endpointLabelMappingFn RequestLabelMappingFn
	consoleColor           bool
	writerLogFn            WriterLogFn
	writerErrorFn          WriterErrorFn
}

// Option for queue system
type Option func(*config)

type WriterLogFn func(log *LogFormatterParams)

type WriterErrorFn func(log *LogFormatterParams) (int, interface{})

// WithLogger set logger function
func WithLogger(logger glog.ILogger) Option {
	return func(cfg *config) {
		cfg.logger = logger.WithField("Gin-Logger", "Gin-Logger")
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

// WithConsoleColor set consoleColor function
func WithConsoleColor(consoleColor bool) Option {
	return func(cfg *config) {
		cfg.consoleColor = consoleColor
	}
}

// WithWriterLogFn set fn WriterLogFn
func WithWriterLogFn(fn WriterLogFn) Option {
	return func(cfg *config) {
		cfg.writerLogFn = fn
	}
}

// WithWriterErrorFn set fn WriterErrorFn
func WithWriterErrorFn(fn WriterErrorFn) Option {
	return func(cfg *config) {
		cfg.writerErrorFn = fn
	}
}
