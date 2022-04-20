package main

import (
	"fmt"
	ginZap "github.com/donetkit/gin-contrib/zap"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	r := gin.New()
	logger, _ := zap.NewProduction()

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginZap.Ginzap(logger, time.RFC3339, false))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	r.Use(ginZap.RecoveryWithZap(logger, false))

	// Example ping request.
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Example when panic happen.
	r.GET("/panic", func(c *gin.Context) {
		panic("An unexpected error happen!")
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
