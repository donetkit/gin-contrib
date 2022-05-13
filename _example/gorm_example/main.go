package main

import (
	"context"
	"github.com/donetkit/gin-contrib-log/glog"
	"github.com/donetkit/gin-contrib/db/gorm"
	"github.com/donetkit/gin-contrib/trace"
)

const (
	service     = "gorm-test"
	environment = "development" // "production" "development"
)

func main() {
	ctx := context.Background()
	log := glog.NewDefaultLogger()
	var traceServer *trace.Server
	tp, err := trace.NewTracerProvider(service, "127.0.0.1", environment, 6831)
	if err == nil {
		jaeger := trace.Jaeger{}
		traceServer = trace.New(service, trace.WithTracerProvider(tp), trace.WithPropagators(jaeger))
	}
	var dns = map[string]string{}
	dns["default"] = "root:test@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms"
	sql := gorm.NewDb(gorm.WithDNS(dns), gorm.WithLogger(log), gorm.WithTracer(traceServer))
	defer func() {
		tp.Shutdown(context.Background())
	}()
	ctx, span := traceServer.Tracer.Start(ctx, "testgorm")
	defer span.End()
	var num []int
	if err := sql.DB().WithContext(ctx).Raw("SELECT holiday FROM calendar where id != ''").Scan(&num).Error; err != nil {
		panic(err)
	}

	var str []string
	if err := sql.DB().WithContext(ctx).Raw("SELECT id FROM school").Scan(&str).Error; err != nil {
		panic(err)
	}

	if err := sql.DB().WithContext(ctx).Raw("SELECT id FROM school_air_conditioner_room").Scan(&str).Error; err != nil {
		panic(err)
	}

	for i := 0; i < 50; i++ {
		if err := sql.DB().WithContext(ctx).Raw("SELECT id FROM school_air_conditioner_room").Scan(&str).Error; err != nil {
			panic(err)
		}
	}

}
