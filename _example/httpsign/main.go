package main

import (
	httpsign2 "github.com/donetkit/gin-contrib/middleware/httpsign"
	crypto2 "github.com/donetkit/gin-contrib/middleware/httpsign/crypto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// Define algorithm
	hmacsha256 := &crypto2.HmacSha256{}
	hmacsha512 := &crypto2.HmacSha512{}
	// Init define secret params
	readKeyID := httpsign2.KeyID("read")
	writeKeyID := httpsign2.KeyID("write")
	secrets := httpsign2.Secrets{
		readKeyID: &httpsign2.Secret{
			Key:       "HMACSHA256-SecretKey",
			Algorithm: hmacsha256, // You could using other algo with interface Crypto
		},
		writeKeyID: &httpsign2.Secret{
			Key:       "HMACSHA512-SecretKey",
			Algorithm: hmacsha512,
		},
	}

	// Init webserve
	r := gin.Default()

	//Create middleware with default rule. Could modify by parse Option func
	auth := httpsign2.NewAuthenticator(secrets)

	r.Use(auth.Authenticated())
	r.GET("/a", a)
	r.POST("/b", b)
	r.POST("/c", c)

	r.Run(":8080")
}

func c(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func b(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func a(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
