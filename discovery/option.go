package discovery

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Option for queue system
type Option func(*Config)

// WithId set id function
func WithId(id string) Option {
	return func(cfg *Config) {
		cfg.Id = id
	}
}

// WithServiceName set serviceName function
func WithServiceName(serviceName string) Option {
	return func(cfg *Config) {
		cfg.ServiceName = serviceName
	}
}

// WithServiceRegisterAddr set serviceRegisterAddr function
func WithServiceRegisterAddr(serviceRegisterAddr string) Option {
	return func(cfg *Config) {
		cfg.ServiceRegisterAddr = serviceRegisterAddr
	}
}

// WithServiceRegisterPort set serviceRegisterPort function
func WithServiceRegisterPort(serviceRegisterPort int) Option {
	return func(cfg *Config) {
		cfg.ServiceRegisterPort = serviceRegisterPort
	}
}

// WithServiceCheckAddr set serviceCheckAddr function
func WithServiceCheckAddr(serviceCheckAddr string) Option {
	return func(cfg *Config) {
		cfg.ServiceCheckAddr = serviceCheckAddr
	}
}

// WithServiceCheckPort set serviceCheckPort function
func WithServiceCheckPort(serviceCheckPort int) Option {
	return func(cfg *Config) {
		cfg.ServiceCheckPort = serviceCheckPort
	}
}

// WithTags set tags function
func WithTags(tags ...string) Option {
	return func(cfg *Config) {
		cfg.Tags = tags
	}
}

// WithIntervalTime set intervalTime function
func WithIntervalTime(intervalTime int) Option {
	return func(cfg *Config) {
		cfg.IntervalTime = intervalTime
	}
}

// WithDeregisterTime set deregisterTime function
func WithDeregisterTime(deregisterTime int) Option {
	return func(cfg *Config) {
		cfg.DeregisterTime = deregisterTime
	}
}

// WithTimeOut set timeOut function
func WithTimeOut(timeOut int) Option {
	return func(cfg *Config) {
		cfg.TimeOut = timeOut
	}
}

// WithCheckHTTP set checkHttp function
func WithCheckHTTP(router *gin.Engine, checkHttp ...string) Option {
	return func(cfg *Config) {
		var checkHttpUrl = "/health"
		if len(checkHttp) > 0 {
			checkHttpUrl = checkHttp[0]
		}
		cfg.CheckHTTP = checkHttpUrl
		cfg.CheckHTTP = fmt.Sprintf("http://%s:%d%s", cfg.ServiceCheckAddr, cfg.ServiceCheckPort, checkHttpUrl)
		router.GET(checkHttpUrl, func(c *gin.Context) { c.String(200, "Healthy") })
	}
}
