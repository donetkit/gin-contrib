package tracer

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	ltrace "go.opentelemetry.io/otel/trace"
)

// New returns middleware that will tracer incoming requests.
// The service parameter should describe the name of the (virtual)
// webserve handling the request.
func New(opts ...Option) *Server {
	cfg := &Server{
		tracerName: "Service",
	}
	for _, opt := range opts {
		opt.apply(cfg)
	}
	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}
	cfg.Tracer = cfg.TracerProvider.Tracer(
		cfg.tracerName,
		ltrace.WithInstrumentationVersion(SemVersion()),
	)
	if cfg.Propagators == nil {
		cfg.Propagators = otel.GetTextMapPropagator()
	}
	return cfg
}

func (s *Server) Stop(ctx context.Context) {
	tp, ok := s.TracerProvider.(*trace.TracerProvider)
	if ok {
		tp.Shutdown(ctx)
	}
}
