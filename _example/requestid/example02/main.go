package main

import (
	"fmt"
	requestid2 "github.com/donetkit/gin-contrib/middleware/requestid"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	r.Use(
		requestid2.New(
			requestid2.WithGenerator(func() string {
				return "test"
			}),
			requestid2.WithCustomHeaderStrKey("your-customer-key"),
		),
	)

	// Example ping request.
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Example / request.
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "id:"+requestid2.Get(c))
	})

	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
