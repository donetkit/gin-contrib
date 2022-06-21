package main

import (
	"github.com/donetkit/contrib-gin/middleware/session"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	store := session.NewCookieStore([]byte("gin-secret"))
	r.Use(session.New("gin-session", store, nil))

	r.GET("/incr", func(c *gin.Context) {
		session := session.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count += 1
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})
	r.Run(":8000")
}
