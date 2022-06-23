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
	r.Use(logger.NewErrorLogger(logger.WithWriterErrorFn(func(c *gin.Context, log *logger.LogFormatterParams) (int, interface{}) {
		//fmt.Println(log)
		return 0, "网络超时, 请重试!"
	})))
	r.Use(logger.New(logger.WithLogger(logs), logger.WithWriterLogFn(func(c *gin.Context, log *logger.LogFormatterParams) {
		//fmt.Println(log)
	})))
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})
	r.GET("/err", func(c *gin.Context) {
		var a = 0
		var b = 1
		fmt.Println(b / a)
		c.String(http.StatusOK, "err "+fmt.Sprint(time.Now().Unix()))
	})
	webserve.New(webserve.WithPort(8080)).AddHandler(r).Run()
}
