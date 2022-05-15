package main

import (
	"github.com/donetkit/gin-contrib/middleware/favicon"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(favicon.New(favicon.WithRoutePaths("/test/favicon.ico")))
	r.Run()
}
