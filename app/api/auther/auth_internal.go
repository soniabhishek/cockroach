package auther

import (
	"fmt"
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/utilities/clients/validator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

func AuthorizeHeader() gin.HandlerFunc {

	jwtKey := config.JWT_SECRET_KEY.Get()

	return func(c *gin.Context) {
		value := c.Request.Header.Get("authorization")
		authHeader := strings.Split(value, " ")
		if len(authHeader) != 2 {
			ShowAuthenticationErrorOverHttp(c, "Auth Failed Invalid Token")
			return
		}
		token, err := jwt.Parse(authHeader[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtKey), nil
		})
		if err != nil || !token.Valid {
			ShowAuthenticationErrorOverHttp(c, "Auth Failed Invalid Access")
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userId", claims["id"])
		} else {
			ShowAuthenticationErrorOverHttp(c, "Invalid Claims")
			return
		}

	}
}
