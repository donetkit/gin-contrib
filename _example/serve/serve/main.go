package main

import (
	"fmt"
	"github.com/donetkit/contrib-log/glog"
	logger2 "github.com/donetkit/gin-contrib/middleware/logger"
	"github.com/donetkit/gin-contrib/server/webserve"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	logs := glog.New()
	r := gin.New()
	r.Use(logger2.New(logger2.WithLogger(logs)))
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})
	webserve.New().AddHandler(r).Run()
}
