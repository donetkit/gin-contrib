package gorm

import (
	"github.com/donetkit/gin-contrib-log/glog"
	"github.com/donetkit/gin-contrib/trace"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"time"
)

type Option func(p *config)

// WithDNS  dsn
func WithDNS(dsn map[string]string) Option {
	return func(p *config) {
		p.dsn = dsn
	}
}

// WithSlowThreshold  slowThreshold
func WithSlowThreshold(slowThreshold time.Duration) Option {
	return func(p *config) {
		p.gormPlugin.slowThreshold = slowThreshold
	}
}

// WithIgnoreRecordNotFoundError  ignoreRecordNotFoundError
func WithIgnoreRecordNotFoundError(ignoreRecordNotFoundError bool) Option {
	return func(p *config) {
		p.gormPlugin.ignoreRecordNotFoundError = ignoreRecordNotFoundError
	}
}

// WithTracer  tracerServer trace.Server
func WithTracer(tracerServer *trace.Server) Option {
	return func(p *config) {
		p.gormPlugin.tracerServer = tracerServer
	}
}

// WithAttributes configures attributes that are used to create a span.
func WithAttributes(attrs ...attribute.KeyValue) Option {
	return func(p *config) {
		p.gormPlugin.attrs = append(p.gormPlugin.attrs, attrs...)
	}
}

// WithDBName configures a db.name attribute.
func WithDBName(name string) Option {
	return func(p *config) {
		p.gormPlugin.attrs = append(p.gormPlugin.attrs, semconv.DBNameKey.String(name))
	}
}

// WithoutQueryVariables configures the db.statement attribute to exclude query variables
func WithoutQueryVariables() Option {
	return func(p *config) {
		p.gormPlugin.excludeQueryVars = true
	}
}

// WithQueryFormatter configures a query formatter
func WithQueryFormatter(queryFormatter func(query string) string) Option {
	return func(p *config) {
		p.gormPlugin.queryFormatter = queryFormatter
	}
}

// WithoutMetrics prevents DBStats metrics from being reported.
func WithoutMetrics() Option {
	return func(p *config) {
		p.gormPlugin.excludeMetrics = true
	}
}

// WithLogger prevents logger.
func WithLogger(logger glog.ILogger) Option {
	return func(p *config) {
		p.gormPlugin.logger = logger
	}
}
