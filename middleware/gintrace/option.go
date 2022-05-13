package gintrace

import (
	"github.com/donetkit/gin-contrib/trace"
)

type config struct {
	tracerServer *trace.Server
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
