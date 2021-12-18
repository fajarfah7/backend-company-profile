package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"project_company_profile/models/user"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type login struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

var identityKey = "username"

func AuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secretKeyForUsers"),
		Timeout:     3 * time.Hour,
		MaxRefresh:  3 * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(user.User); ok {
				return jwt.MapClaims{
					// "id":           v.ID,
					identityKey: v.Username,
					// "name":         v.Name,
					// "username":     v.Username,
					// "email":        v.Email,
					// "address":      v.Address,
					// "phone_number": v.PhoneNumber,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			// get or create cart on here
			return &user.User{
				Username: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			requestVars := map[string]string{"username": username, "password": password}
			requestJSON, err := json.Marshal(requestVars)
			if err != nil {
				return "", err
			}

			payload := bytes.NewBufferString(string(requestJSON))

			req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:4000/user/get-by-username-and-password", payload)
			if err != nil {
				return "", err
			}
			req.Header.Add("Content-Type", "application/json")

			client := &http.Client{}

			res, err := client.Do(req)
			if err != nil {
				return "", err
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				return "", err
			}

			if res.StatusCode != http.StatusOK {
				type BodyMessage struct {
					Message string `json:"message" form:"message"`
				}

				var bodyMessage BodyMessage
				err := json.Unmarshal(body, &bodyMessage)
				if err != nil {
					return "", err
				}

				return "", errors.New(bodyMessage.Message)
			}

			var user user.User
			err = json.Unmarshal(body, &user)
			if err != nil {
				return "", err
			}

			return user, nil
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"status":   code,
				"messages": []map[string]string{{"api": message}},
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie:jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}
