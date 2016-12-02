package auther

import (
	"fmt"
	"github.com/crowdflux/angel/app/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

func AuthorizeHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.Request.Header.Get("authorization")
		authHeader := strings.Split(value, " ")

		token, err := jwt.Parse(authHeader[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.JWT_SECRET_KEY.Get()), nil
		})
		if err != nil || !token.Valid {
			fmt.Println(err)
			c.Header("authenication error", "Invalid Access")
			c.AbortWithStatus(401)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userId", claims["id"])
		} else {
			c.Header("authenication error", "Invalid Claims")
			c.AbortWithStatus(401)
			return
		}

	}
}
