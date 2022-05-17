package main

import (
	"github.com/donetkit/contrib-gin/middleware/pprof"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	pprof.Register(router)
	router.Run(":8080")
}
