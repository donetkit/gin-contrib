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
		p.gormConfig.slowThreshold = slowThreshold
	}
}

// WithIgnoreRecordNotFoundError  ignoreRecordNotFoundError
func WithIgnoreRecordNotFoundError(ignoreRecordNotFoundError bool) Option {
	return func(p *config) {
		p.gormConfig.ignoreRecordNotFoundError = ignoreRecordNotFoundError
	}
}

// WithTracer  tracerServer trace.Server
func WithTracer(tracerServer *trace.Server) Option {
	return func(p *config) {
		p.gormConfig.tracerServer = tracerServer
	}
}

// WithAttributes configures attributes that are used to create a span.
func WithAttributes(attrs ...attribute.KeyValue) Option {
	return func(p *config) {
		p.gormConfig.attrs = append(p.gormConfig.attrs, attrs...)
	}
}

// WithDBName configures a db.name attribute.
func WithDBName(name string) Option {
	return func(p *config) {
		p.gormConfig.attrs = append(p.gormConfig.attrs, semconv.DBNameKey.String(name))
	}
}

// WithoutQueryVariables configures the db.statement attribute to exclude query variables
func WithoutQueryVariables() Option {
	return func(p *config) {
		p.gormConfig.excludeQueryVars = true
	}
}

// WithQueryFormatter configures a query formatter
func WithQueryFormatter(queryFormatter func(query string) string) Option {
	return func(p *config) {
		p.gormConfig.queryFormatter = queryFormatter
	}
}

// WithoutMetrics prevents DBStats metrics from being reported.
func WithoutMetrics() Option {
	return func(p *config) {
		p.gormConfig.excludeMetrics = true
	}
}

// WithLogger prevents logger.
func WithLogger(logger glog.ILogger) Option {
	return func(p *config) {
		p.gormConfig.logger = logger
	}
}

// WithDefaultStringSize  defaultStringSize
func WithDefaultStringSize(defaultStringSize uint) Option {
	return func(p *config) {
		p.gormConfig.defaultStringSize = defaultStringSize
	}
}

// WithDisableDatetimePrecision  disableDatetimePrecision
func WithDisableDatetimePrecision(disableDatetimePrecision bool) Option {
	return func(p *config) {
		p.gormConfig.disableDatetimePrecision = disableDatetimePrecision
	}
}

// WithDontSupportRenameIndex  dontSupportRenameIndex
func WithDontSupportRenameIndex(dontSupportRenameIndex bool) Option {
	return func(p *config) {
		p.gormConfig.dontSupportRenameIndex = dontSupportRenameIndex
	}
}

// WithDontSupportRenameColumn  dontSupportRenameColumn
func WithDontSupportRenameColumn(dontSupportRenameColumn bool) Option {
	return func(p *config) {
		p.gormConfig.dontSupportRenameColumn = dontSupportRenameColumn
	}
}

// WithSkipInitializeWithVersion  skipInitializeWithVersion
func WithSkipInitializeWithVersion(skipInitializeWithVersion bool) Option {
	return func(p *config) {
		p.gormConfig.skipInitializeWithVersion = skipInitializeWithVersion
	}
}
