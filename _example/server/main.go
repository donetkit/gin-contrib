package main

import (
	"fmt"
	"github.com/donetkit/contrib-log/glog"
	"github.com/donetkit/gin-contrib/discovery"
	"github.com/donetkit/gin-contrib/discovery/consul"
	logger2 "github.com/donetkit/gin-contrib/middleware/logger"
	"github.com/donetkit/gin-contrib/server"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	logs := glog.New()
	consulClient, _ := consul.New(discovery.WithServiceRegisterAddr("127.0.0.1"))
	r := gin.New()
	r.Use(logger2.New(logger2.WithLogger(logs)))
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})
	appServe, err := server.New()
	if err != nil {
		panic(err)
	}
	appServe.AddDiscovery(consulClient).AddHandler(r).Run()
}
