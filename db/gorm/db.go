package gorm

import (
	"github.com/donetkit/gin-contrib-log/glog"
	"github.com/donetkit/gin-contrib/trace"
	"github.com/prometheus/common/log"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type config struct {
	dsn        map[string]string
	gormPlugin *gormPlugin
}

type gormPlugin struct {
	ignoreRecordNotFoundError bool
	slowThreshold             time.Duration
	logger                    glog.ILogger
	tracerServer              *trace.Server
	attrs                     []attribute.KeyValue
	excludeQueryVars          bool
	excludeMetrics            bool
	queryFormatter            func(query string) string
}

func NewDb(opts ...Option) map[string]*gorm.DB {
	p := &config{}
	p.gormPlugin = &gormPlugin{}
	for _, opt := range opts {
		opt(p)
	}
	gormConfig := &gorm.Config{}
	if p.gormPlugin.logger != nil {
		gormConfig.Logger = &LogSql{Logger: p.gormPlugin.logger, config: p.gormPlugin}
	}
	var dbs = map[string]*gorm.DB{}
	for key, val := range p.dsn {
		db, err := gorm.Open(mysql.New(mysql.Config{
			DSN:                       val,
			DefaultStringSize:         256,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		}), gormConfig)
		if err != nil && p.gormPlugin.logger != nil {
			log.Error(err.Error())
		}
		if p.gormPlugin.tracerServer != nil {
			if err := db.Use(NewPlugin(p.gormPlugin)); err != nil && p.gormPlugin.logger != nil {
				log.Error(err.Error())
			}
		}
		dbs[key] = db
	}
	return dbs
}
