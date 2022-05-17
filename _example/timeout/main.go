package main

import (
	"github.com/donetkit/contrib-gin/middleware/timeout"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func emptySuccessResponse(c *gin.Context) {
	time.Sleep(5 * time.Second)
	c.String(http.StatusOK, "OK")
}

func main() {
	r := gin.New()

	r.GET("/", timeout.New(
		timeout.WithTimeout(3*time.Second),
		timeout.WithHandler(emptySuccessResponse),
	))

	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
