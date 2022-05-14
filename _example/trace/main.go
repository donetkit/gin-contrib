package main

import (
	"context"
	"github.com/donetkit/gin-contrib/middleware/gintrace"
	"github.com/donetkit/gin-contrib/tracer"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	service     = "gin-gonic-development-server"
	environment = "development" // "production" "development"
)

func main() {
	r := gin.New()

	tp, err := tracer.NewTracerProvider(service, "127.0.0.1", environment, 6831)
	if err == nil {
		jaeger := tracer.Jaeger{}
		traceServer := tracer.New(service, tracer.WithTracerProvider(tp), tracer.WithPropagators(jaeger))
		r.Use(gintrace.New(service, gintrace.WithTracer(traceServer)))
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
