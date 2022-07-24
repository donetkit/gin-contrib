package main

import (
	"context"
	"github.com/donetkit/contrib-gin/middleware/gintrace"
	"github.com/donetkit/contrib/tracer"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	service     = "gin-gonic-development-webserve"
	environment = "development" // "production" "development"
)

func main() {
	r := gin.New()

	tp, err := tracer.NewTracerProvider(service, "127.0.0.1", environment, 6831, nil)
	if err == nil {
		jaeger := tracer.Jaeger{}
		traceServer := tracer.New(tracer.WithName(service), tracer.WithProvider(tp), tracer.WithPropagators(jaeger))
		r.Use(gintrace.New(gintrace.WithName(service), gintrace.WithTracer(traceServer), gintrace.WithWriterTraceId(), gintrace.WithWriterSpanId()))
		defer func() {
			tp.Shutdown(context.Background())
		}()
	}

	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	})
	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
