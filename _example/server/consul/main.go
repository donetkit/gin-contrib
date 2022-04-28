package main

import (
	"fmt"
	"github.com/donetkit/gin-contrib/discovery/consul"
	"github.com/donetkit/gin-contrib/server"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	r := gin.New()
	consulClient, _ := consul.New()
	// Example ping request.
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})
	appServe, err := server.New(server.WithRouter(r), server.WithConsulClient(consulClient))
	if err != nil {
		panic(err)
	}
	appServe.Start()
	appServe.AwaitSignal()
}
