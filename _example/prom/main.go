package main

import (
	"github.com/donetkit/contrib-gin/middleware/prom"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()

	router.Use(prom.New(prom.WithNamespace("service"), prom.WithPromHandler(router)))

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
