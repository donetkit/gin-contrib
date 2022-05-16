package cache

import (
	"github.com/donetkit/gin-contrib-log/glog"
	"github.com/donetkit/gin-contrib/utils/cache"
)

type option struct {
	cache  cache.IShortCache
	logger glog.ILogger
}

type Option func(*option)

func WithLogger(logger glog.ILogger) Option {
	return func(o *option) {
		o.logger = logger
	}
}

// WithCache  cache
func WithCache(cache cache.IShortCache) Option {
	return func(o *option) {
		o.cache = cache
	}
}
