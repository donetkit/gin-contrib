package main

import (
	"fmt"
	"github.com/donetkit/gin-contrib/discovery/consul"
	server2 "github.com/donetkit/gin-contrib/server/webserve"
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
	server2.New(server2.WithHandler(r)).AddDiscovery(consulClient).Run()

}
