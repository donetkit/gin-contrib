package main

import (
	"github.com/donetkit/contrib-gin/middleware/favicon"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(favicon.New(favicon.WithRoutePaths("/test/favicon.ico")))
	r.Run()
}
