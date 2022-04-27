package main

import (
	"context"
	"github.com/donetkit/gin-contrib/trace"
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
	tp, err := trace.NewTracerProvider(service, "192.168.5.110", environment, 6831)
	if err == nil {
		jaeger := trace.Jaeger{}
		r.Use(trace.Middleware(service, trace.WithTracerProvider(tp), trace.WithPropagators(jaeger)))
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
