package main

import (
	"fmt"
	"github.com/donetkit/contrib/discovery"
	"github.com/donetkit/contrib/discovery/consul"
	"github.com/donetkit/contrib/server/webserve"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	r := gin.New()

	client, _ := consul.New(
		discovery.WithRegisterAddr("127.0.0.1"),
		discovery.WithRegisterPort(8500),
		discovery.WithEnableHealthyStatus(),
		discovery.WithCheckHTTP(func(response *discovery.CheckResponse) {
			r.GET(response.Url, func(c *gin.Context) { c.String(200, response.Result()) })
		}))
	// Example ping request.
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})
	webserve.New(webserve.WithHandler(r)).AddDiscovery(client).Run()

}
