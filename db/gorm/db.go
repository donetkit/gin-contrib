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
	gormConfig *gormConfig
}

type gormConfig struct {
	ignoreRecordNotFoundError bool
	slowThreshold             time.Duration
	logger                    glog.ILogger
	tracerServer              *trace.Server
	attrs                     []attribute.KeyValue
	excludeQueryVars          bool
	excludeMetrics            bool
	queryFormatter            func(query string) string
	defaultStringSize         uint // string 类型字段的默认长度
	disableDatetimePrecision  bool // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	dontSupportRenameIndex    bool // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
	dontSupportRenameColumn   bool // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
	skipInitializeWithVersion bool // 根据当前 MySQL 版本自动配置

	connMaxIdleTime time.Duration // 设置了连接可复用的最大时间
	maxOpenCons     int           // 设置打开数据库连接的最大数量
	maxIdleCons     int           // 设置空闲连接池中连接的最大数量
}

func NewDb(opts ...Option) map[string]*gorm.DB {
	p := &config{}
	p.gormConfig = &gormConfig{
		defaultStringSize:         256,
		disableDatetimePrecision:  true,
		dontSupportRenameIndex:    true,
		dontSupportRenameColumn:   true,
		skipInitializeWithVersion: false,
		connMaxIdleTime:           time.Second * 1800,
		maxOpenCons:               100,
		maxIdleCons:               20,
	}
	for _, opt := range opts {
		opt(p)
	}
	gormConfig := &gorm.Config{}
	if p.gormConfig.logger != nil {
		gormConfig.Logger = &LogSql{Logger: p.gormConfig.logger, config: p.gormConfig}
	}
	var dbs = map[string]*gorm.DB{}
	for key, val := range p.dsn {
		db, err := gorm.Open(mysql.New(mysql.Config{
			DSN:                       val,
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}), gormConfig)
		if err != nil && p.gormConfig.logger != nil {
			log.Error(err.Error())
		}
		if p.gormConfig.tracerServer != nil {
			if err := db.Use(NewPlugin(p.gormConfig)); err != nil && p.gormConfig.logger != nil {
				log.Error(err.Error())
			}
		}
		sdb, err := db.DB()
		if err == nil {
			sdb.SetConnMaxIdleTime(p.gormConfig.connMaxIdleTime) //最大生存时间(s) 30 SetConnMaxLifetime 设置了连接可复用的最大时间。
			sdb.SetMaxOpenConns(p.gormConfig.maxOpenCons)        // SetMaxOpenConns 设置打开数据库连接的最大数量。
			sdb.SetMaxIdleConns(p.gormConfig.maxIdleCons)        //最大连接数 1000 SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
		}
		dbs[key] = db
	}
	return dbs
}
