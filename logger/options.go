package logger

import "github.com/donetkit/gin-contrib/utils/glog"

// Config defines the config for logger middleware
type config struct {
	logger glog.ILogger
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
