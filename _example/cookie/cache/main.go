package main

import (
	"github.com/donetkit/contrib-gin/middleware/session"
	"github.com/donetkit/contrib/db/memory"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//cacheClient := redis.New(redis.WithAddr("192.168.0.3"), redis.WithPort(6379), redis.WithPassword("test")).WithDB(2)
	cacheClient := memory.New().WithDB(2)
	store, _ := session.NewStore(cacheClient, []byte("gin-secret"))
	r.Use(session.New("gin-session-cache", store, nil))

	r.GET("/incr", func(c *gin.Context) {
		session := session.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count1, ok := v.(float64)
			if ok {
				count = int(count1) + 1
			}
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})
	r.Run(":8000")
}
