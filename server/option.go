package server

import (
	"github.com/donetkit/gin-contrib/discovery/consul"
	"github.com/donetkit/gin-contrib/utils/glog"
	"github.com/gin-gonic/gin"
	"time"
)

// Option for queue system
type Option func(*config)

// WithServiceName set serviceName function
func WithServiceName(serviceName string) Option {
	return func(cfg *config) {
		cfg.serviceName = serviceName
	}
}

// WithHost set host function
func WithHost(host string) Option {
	return func(cfg *config) {
		cfg.host = host
	}
}

// WithPort set port function
func WithPort(port int) Option {
	return func(cfg *config) {
		cfg.port = port
	}
}

// WithRouter set router function
func WithRouter(router *gin.Engine) Option {
	return func(cfg *config) {
		cfg.router = router
	}
}

// WithHttpServer set httpServer function
//func WithHttpServer(httpServer http.Server) Option {
//	return func(cfg *config) {
//		cfg.httpServer = httpServer
//	}
//}

// WithConsulClient set consulClient function
func WithConsulClient(consulClient *consul.Client) Option {
	return func(cfg *config) {
		cfg.consulClient = consulClient
	}
}

// WithReadTimeout set readTimeout function
func WithReadTimeout(readTimeout time.Duration) Option {
	return func(cfg *config) {
		cfg.readTimeout = readTimeout
	}
}

// WithWriterTimeout set writerTimeout function
func WithWriterTimeout(writerTimeout time.Duration) Option {
	return func(cfg *config) {
		cfg.writerTimeout = writerTimeout
	}
}

// WithMaxHeaderBytes set maxHeaderBytes function
func WithMaxHeaderBytes(maxHeaderBytes int) Option {
	return func(cfg *config) {
		cfg.maxHeaderBytes = maxHeaderBytes
	}
}

// WithLogger set logger function
func WithLogger(logger glog.ILogger) Option {
	return func(cfg *config) {
		cfg.logger = logger
	}
}
