package auther

import (
	"fmt"
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/utilities/clients/validator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

func AuthorizeHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.Request.Header.Get("authorization")
		authHeader := strings.Split(value, " ")
		if len(authHeader) != 2 {
			validator.ShowErrorResponseOverHttp(c, errors.New("Auth Failed"))
			c.Abort()
			return
		}
		token, err := jwt.Parse(authHeader[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.JWT_SECRET_KEY.Get()), nil
		})
		if err != nil || !token.Valid {
			c.Header("authenication error", "Invalid Access")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userId", claims["id"])
		} else {
			c.Header("authenication error", "Invalid Claims")
			c.Abort()
			return
		}

	}
}
