package main

import (
	"fmt"
	logger2 "github.com/donetkit/gin-contrib/middleware/logger"
	"github.com/donetkit/gin-contrib/server"
	"github.com/donetkit/gin-contrib/utils/glog"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	logs := glog.NewDefaultLogger()
	r := gin.New()
	r.Use(logger2.New(logger2.WithLogger(logs)))
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})
	appServe, err := server.New(server.WithRouter(r))
	if err != nil {
		panic(err)
	}
	appServe.Start()
}
