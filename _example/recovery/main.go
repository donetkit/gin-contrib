package main

import (
	"fmt"
	"github.com/donetkit/contrib-gin/middleware/recovery"
	"github.com/donetkit/contrib-log/glog"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	logger := glog.New()
	r.Use(recovery.New(logger, true))
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
