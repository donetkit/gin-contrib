package trace

import (
	"go.opentelemetry.io/otel"
	ltrace "go.opentelemetry.io/otel/trace"
)

// New returns middleware that will trace incoming requests.
// The service parameter should describe the name of the (virtual)
// server handling the request.
func New(tracerName string, opts ...Option) *Server {
	cfg := &Server{}
	for _, opt := range opts {
		opt.apply(cfg)
	}
	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}
	cfg.Tracer = cfg.TracerProvider.Tracer(
		tracerName,
		ltrace.WithInstrumentationVersion(SemVersion()),
	)
	if cfg.Propagators == nil {
		cfg.Propagators = otel.GetTextMapPropagator()
	}
	return cfg
}
