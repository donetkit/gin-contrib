package main

import (
	"context"
	"github.com/donetkit/contrib-gin/gorm_mysql"
	"github.com/donetkit/contrib-log/glog"
	"github.com/donetkit/contrib/tracer"
)

const (
	service     = "gorm-test"
	environment = "development" // "production" "development"
)

func main() {
	ctx := context.Background()
	log := glog.New()
	var traceServer *tracer.Server
	tp, err := tracer.NewTracerProvider(service, "127.0.0.1", environment, 6831, tracer.NewFallbackSampler(1.0))
	if err == nil {
		jaeger := tracer.Jaeger{}
		traceServer = tracer.New(tracer.WithName(service), tracer.WithProvider(tp), tracer.WithPropagators(jaeger))
	}
	var dns = map[string]string{}
	dns["default"] = "root:test@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms"
	sql := gorm_mysql.NewDb(gorm_mysql.WithDNS(dns), gorm_mysql.WithLogger(log), gorm_mysql.WithTracer(traceServer))
	defer func() {
		tp.Shutdown(context.Background())
	}()
	ctx, span := traceServer.Tracer.Start(ctx, "testgorm")
	defer span.End()
	var num []int
	if err := sql.DB().WithContext(ctx).Raw("SELECT id FROM test").Scan(&num).Error; err != nil {
		log.Error(err)
	}

}
