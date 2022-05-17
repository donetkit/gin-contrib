package main

import (
	"fmt"
	"github.com/donetkit/contrib-gin/middleware/gzip"
	"github.com/donetkit/contrib-gin/middleware/logger"
	"github.com/donetkit/contrib-log/glog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	log := glog.New()
	logger.SetGinDefaultWriter(log)
	r := gin.New()
	// LoggerWithFormatter 中间件会将日志写入 gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout

	//gin.DefaultWriter =
	r.Use(gzip.Gzip(gzip.DefaultCompression), logger.New(logger.WithLogger(log)))

	//
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})
	r.GET("/ping2", func(c *gin.Context) {
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})
	r.GET("/ping3", func(c *gin.Context) {
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})
	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
