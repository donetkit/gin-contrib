package main

import (
	"fmt"
	"github.com/donetkit/contrib-gin/middleware/logger"
	"github.com/donetkit/contrib-log/glog"
	"github.com/donetkit/contrib/server/webserve"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	logs := glog.New()
	r := gin.New()
	r.Use(logger.New(logger.WithLogger(logs)))
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})
	webserve.New().AddHandler(r).Run()
}
