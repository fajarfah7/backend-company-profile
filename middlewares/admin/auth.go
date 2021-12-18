package admin

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"

// User is struct that will be used to bind admin data
type User struct {
	Username string
	Name     string
}

// AuthMiddleware place where admin data stored, logic, and some config here
func AuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("thisismysecretkey123"),
		Timeout:     3 * time.Hour,
		MaxRefresh:  3 * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
					"name":      v.Name,
					// identityKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			// if PayloadFunc return jwt.MapClaims then on here will give
			// type map[string]interface{}, so we need to cast it first
			// in order can be used
			// user := claims[identityKey].(map[string]interface{})
			return &User{
				Username: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if userID == "admin" && password == "securepassword" {
				return &User{
					Username: userID,
					Name:     "Fajar Fahrurozi",
				}, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && v.Username == "admin" {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}
