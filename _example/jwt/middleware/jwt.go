package middleware

import (
	"github.com/donetkit/gin-contrib/jwt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

var identityKey = "id"

func NewJWT() *jwt.GinJWTMiddleware {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		TimeFunc:    time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	authMiddleware.PayloadFunc = func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*User); ok {
			return jwt.MapClaims{
				identityKey: v.UserName,
			}
		}
		return jwt.MapClaims{}
	}
	authMiddleware.IdentityHandler = func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		return &User{
			UserName: claims[identityKey].(string),
		}
	}
	authMiddleware.Authenticator = func(c *gin.Context) (interface{}, error) {
		var loginVals login
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		userID := loginVals.Username
		password := loginVals.Password

		if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
			return &User{
				UserName:  userID,
				LastName:  "Bo-Yi",
				FirstName: "Wu",
			}, nil
		}

		return nil, jwt.ErrFailedAuthentication
	}
	authMiddleware.Authorizator = func(data interface{}, c *gin.Context) bool {
		if v, ok := data.(*User); ok && v.UserName == "admin" {
			return true
		}
		return false
	}
	authMiddleware.Unauthorized = func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}
	return authMiddleware
}
