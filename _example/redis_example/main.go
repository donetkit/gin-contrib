package main

import (
	"context"
	"github.com/donetkit/gin-contrib-log/glog"
	redisRedis "github.com/donetkit/gin-contrib/db/redis"
	"github.com/donetkit/gin-contrib/tracer"
	"github.com/donetkit/gin-contrib/utils/cache"
	"time"
)

const (
	service     = "redis-test"
	environment = "development" // "production" "development"
)

func main() {
	ctx := context.Background()

	log := glog.New()
	var traceServer *tracer.Server
	tp, err := tracer.NewTracerProvider(service, "127.0.0.1", environment, 6831)
	if err == nil {
		jaeger := tracer.Jaeger{}
		traceServer = tracer.New(service, tracer.WithTracerProvider(tp), tracer.WithPropagators(jaeger))
	}

	rdb := redisRedis.New(redisRedis.WithLogger(log), redisRedis.WithTracer(traceServer), redisRedis.WithAddr("127.0.0.1"), redisRedis.WithPort(6379), redisRedis.WithPassword("test"), redisRedis.WithDB(0))
	if err := redisCommands(ctx, traceServer, rdb); err != nil {
		log.Error(err.Error())
		return
	}
	log.Info("111111111111111")
	time.Sleep(time.Second * 31)
	//fmt.Println("tracer", otelplay.TraceURL(span))
}
func redisCommands(ctx context.Context, traceServer *tracer.Server, rdb cache.ICache) error {
	ctx, span := traceServer.Tracer.Start(ctx, "11111111111111111111111")
	defer span.End()
	if err := rdb.Set(0, "foo", "bar", 0, ctx); err != nil {
		return err
	}
	rdb.Get(0, "foo", ctx)

	return nil
}
