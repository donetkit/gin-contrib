package main

import (
	"context"
	"fmt"
	"github.com/donetkit/gin-contrib-log/glog"
	"github.com/donetkit/gin-contrib/db/gorm"
	"github.com/donetkit/gin-contrib/trace"
)

const (
	service     = "app_or_package_name"
	environment = "development" // "production" "development"
)

func main() {
	ctx := context.Background()
	log := glog.NewDefaultLogger()
	var traceServer *trace.Server
	tp, err := trace.NewTracerProvider(service, "192.168.5.110", environment, 6831)
	if err == nil {
		jaeger := trace.Jaeger{}
		traceServer = trace.New(service, trace.WithTracerProvider(tp), trace.WithPropagators(jaeger))
	}

	var dns = map[string]string{}
	dns["default"] = "root:zxw123456@tcp(192.168.5.110:3306)/go_red_sentinel?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms"
	dbs := gorm.NewDb(gorm.WithDNS(dns), gorm.WithLogger(log), gorm.WithTracer(traceServer))

	defer func() {
		tp.Shutdown(context.Background())
	}()
	ctx, span := traceServer.Tracer.Start(ctx, "calendar")
	defer span.End()

	for _, db := range dbs {
		var num []int
		if err := db.WithContext(ctx).Raw("SELECT holiday FROM  calendar").Scan(&num).Error; err != nil {
			panic(err)
		}
		fmt.Println(num)
	}

}
