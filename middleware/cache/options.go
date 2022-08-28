package cache

import (
	"github.com/donetkit/contrib-log/glog"
	"github.com/donetkit/contrib/utils/cache"
)

type option struct {
	cache  cache.ICache
	logger glog.ILoggerEntry
}

type Option func(*option)

func WithLogger(logger glog.ILogger) Option {
	return func(o *option) {
		o.logger = logger.WithField("Cache", "Cache")
	}
}

// WithCache  cache
func WithCache(cache cache.ICache) Option {
	return func(o *option) {
		o.cache = cache
	}
}
