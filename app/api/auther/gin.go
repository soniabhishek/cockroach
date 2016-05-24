package auther

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

func GinAuther() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("called")

		clId := c.Request.Header.Get("client_id")
		if clId == "" {
			c.Header("authenication error", "client_id required")
			c.AbortWithStatus(401)
			return
		}

		clSec := c.Request.Header.Get("client_secret")
		if clSec == "" {
			c.Header("authenication error", "client_secret required")
			c.AbortWithStatus(401)
			return

		}

		id, err := uuid.FromCEnc(clId)
		if err != nil {
			c.Header("authenication error", "client_id invalid")
			c.AbortWithStatus(401)
			return

		}

		if StdProdAuther.Check(id, clSec) {

			//Let the request go forward, set a user_id param also
			c.Set("user_id", id.String())
		} else {
			c.Header("authenication error", "unauthorized")
			c.AbortWithStatus(401)
		}

	}
}
