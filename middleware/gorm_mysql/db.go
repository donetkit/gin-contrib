package gorm_mysql

import (
	"github.com/donetkit/contrib-log/glog"
	ggrom "github.com/donetkit/contrib/db/gorm"
	"github.com/donetkit/contrib/tracer"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Client struct {
	Client map[string]*gorm.DB
	config *config
}

func (c *Client) DB(dbName ...string) (conn *gorm.DB) {
	if c.Client == nil {
		return nil
	}
	if len(dbName) <= 0 {
		return c.Client["default"]
	}
	return c.Client[dbName[0]]
}

type config struct {
	dsn           map[string]string
	slowThreshold time.Duration
	logger        glog.ILogger
	tracerServer  *tracer.Server

	sqlConfig                 *sqlConfig
	ignoreRecordNotFoundError bool

	attrs            []attribute.KeyValue
	excludeQueryVars bool
	excludeMetrics   bool
	queryFormatter   func(query string) string

	defaultStringSize         uint // string 类型字段的默认长度
	disableDatetimePrecision  bool // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	dontSupportRenameIndex    bool // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
	dontSupportRenameColumn   bool // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
	skipInitializeWithVersion bool // 根据当前 MySQL 版本自动配置

	connMaxIdleTime time.Duration // 设置了连接可复用的最大时间
	maxOpenCons     int           // 设置打开数据库连接的最大数量
	maxIdleCons     int           // 设置空闲连接池中连接的最大数量
}

type sqlConfig struct {
}

func NewDb(opts ...Option) *Client {
	c := &config{
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
		opt(c)
	}
	cfg := &gorm.Config{}
	if c.logger != nil {
		cfg.Logger = &ggrom.LogSql{Logger: c.logger, IgnoreRecordNotFoundError: c.ignoreRecordNotFoundError, SlowThreshold: c.slowThreshold}
	}
	var dbs = map[string]*gorm.DB{}
	for key, val := range c.dsn {
		db, err := gorm.Open(mysql.New(mysql.Config{
			DSN:                       val,
			DefaultStringSize:         c.defaultStringSize,         // string 类型字段的默认长度
			DisableDatetimePrecision:  c.disableDatetimePrecision,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    c.dontSupportRenameIndex,    // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   c.dontSupportRenameColumn,   // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: c.skipInitializeWithVersion, // 根据当前 MySQL 版本自动配置
		}), cfg)
		if err != nil && c.logger != nil {
			c.logger.Error(err)
		}
		if c.tracerServer != nil {
			gcfg := &ggrom.Config{
				Logger:           c.logger,
				TracerServer:     c.tracerServer,
				Attrs:            c.attrs,
				ExcludeMetrics:   c.excludeMetrics,
				ExcludeQueryVars: c.excludeQueryVars,
				QueryFormatter:   c.queryFormatter,
			}
			if err := db.Use(ggrom.NewPlugin(gcfg)); err != nil && c.logger != nil {
				c.logger.Error(err)
			}
		}
		sdb, err := db.DB()
		if err == nil {
			sdb.SetConnMaxIdleTime(c.connMaxIdleTime) //最大生存时间(s) 30 SetConnMaxLifetime 设置了连接可复用的最大时间。
			sdb.SetMaxOpenConns(c.maxOpenCons)        // SetMaxOpenConns 设置打开数据库连接的最大数量。
			sdb.SetMaxIdleConns(c.maxIdleCons)        //最大连接数 1000 SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
		}

		dbs[key] = db
	}
	return &Client{Client: dbs, config: c}
}
