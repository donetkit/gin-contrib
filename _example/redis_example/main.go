package main

import (
	"context"
	"github.com/donetkit/gin-contrib-log/glog"
	redisRedis "github.com/donetkit/gin-contrib/db/redis"
	"github.com/donetkit/gin-contrib/trace"
	"time"
)

const (
	service     = "redis-test"
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

	rdb := redisRedis.New(redisRedis.WithLogger(log), redisRedis.WithTracer(traceServer), redisRedis.WithAddr("127.0.0.1"), redisRedis.WithPort(6379), redisRedis.WithPassword("test"), redisRedis.WithDB(0))

	if err := redisCommands(ctx, traceServer, rdb); err != nil {
		log.Error(err.Error())
		return
	}
	log.Info("111111111111111")
	time.Sleep(time.Second * 31)
	//fmt.Println("trace", otelplay.TraceURL(span))
}
func redisCommands(ctx context.Context, traceServer *trace.Server, rdb *redisRedis.Client) error {
	ctx, span := traceServer.Tracer.Start(ctx, "11111111111111111111111")
	defer span.End()
	if err := rdb.Set(0, "foo", "bar", 0, ctx); err != nil {
		return err
	}
	rdb.Get(0, "foo", ctx)

	return nil
}
