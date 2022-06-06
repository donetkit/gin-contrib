package gorm_mysql

import (
	"github.com/donetkit/contrib-log/glog"
	"github.com/donetkit/contrib/tracer"
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
		p.slowThreshold = slowThreshold
	}
}

// WithIgnoreRecordNotFoundError  ignoreRecordNotFoundError
func WithIgnoreRecordNotFoundError(ignoreRecordNotFoundError bool) Option {
	return func(p *config) {
		p.ignoreRecordNotFoundError = ignoreRecordNotFoundError
	}
}

// WithTracer  tracerServer tracer.Server
func WithTracer(tracerServer *tracer.Server) Option {
	return func(p *config) {
		p.tracerServer = tracerServer
	}
}

// WithAttributes configures attributes that are used to create a span.
func WithAttributes(attrs ...attribute.KeyValue) Option {
	return func(p *config) {
		p.attrs = append(p.attrs, attrs...)
	}
}

// WithDBName configures a db.name attribute.
func WithDBName(name string) Option {
	return func(p *config) {
		p.attrs = append(p.attrs, semconv.DBNameKey.String(name))
	}
}

// WithoutQueryVariables configures the db.statement attribute to exclude query variables
func WithoutQueryVariables() Option {
	return func(p *config) {
		p.excludeQueryVars = true
	}
}

// WithQueryFormatter configures a query formatter
func WithQueryFormatter(queryFormatter func(query string) string) Option {
	return func(p *config) {
		p.queryFormatter = queryFormatter
	}
}

// WithoutMetrics prevents DBStats metrics from being reported.
func WithoutMetrics() Option {
	return func(p *config) {
		p.excludeMetrics = true
	}
}

// WithLogger prevents logger.
func WithLogger(logger glog.ILogger) Option {
	return func(p *config) {
		p.logger = logger
	}
}

// WithDefaultStringSize  defaultStringSize
func WithDefaultStringSize(defaultStringSize uint) Option {
	return func(p *config) {
		p.defaultStringSize = defaultStringSize
	}
}

// WithDisableDatetimePrecision  disableDatetimePrecision
func WithDisableDatetimePrecision(disableDatetimePrecision bool) Option {
	return func(p *config) {
		p.disableDatetimePrecision = disableDatetimePrecision
	}
}

// WithDontSupportRenameIndex  dontSupportRenameIndex
func WithDontSupportRenameIndex(dontSupportRenameIndex bool) Option {
	return func(p *config) {
		p.dontSupportRenameIndex = dontSupportRenameIndex
	}
}

// WithDontSupportRenameColumn  dontSupportRenameColumn
func WithDontSupportRenameColumn(dontSupportRenameColumn bool) Option {
	return func(p *config) {
		p.dontSupportRenameColumn = dontSupportRenameColumn
	}
}

// WithSkipInitializeWithVersion  skipInitializeWithVersion
func WithSkipInitializeWithVersion(skipInitializeWithVersion bool) Option {
	return func(p *config) {
		p.skipInitializeWithVersion = skipInitializeWithVersion
	}
}

// WithConnMaxIdleTime  connMaxIdleTime
func WithConnMaxIdleTime(connMaxIdleTime time.Duration) Option {
	return func(p *config) {
		p.connMaxIdleTime = connMaxIdleTime
	}
}

// WithMaxOpenCons  maxOpenCons
func WithMaxOpenCons(maxOpenCons int) Option {
	return func(p *config) {
		p.maxOpenCons = maxOpenCons
	}
}

// WithMaxIdleCons  maxIdleCons
func WithMaxIdleCons(maxIdleCons int) Option {
	return func(p *config) {
		p.maxIdleCons = maxIdleCons
	}
}
