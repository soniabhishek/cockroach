package auther

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

func GinAuther() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get Client ID
		clId := c.Request.Header.Get("client_id")
		if clId == "" {
			c.Header("authenication error", "client_id required")
			c.AbortWithStatus(401)
			return
		}

		// Get Client Secret
		clSec := c.Request.Header.Get("client_secret")
		if clSec == "" {
			c.Header("authenication error", "client_secret required")
			c.AbortWithStatus(401)
			return

		}

		// Get UUID from ClientID (client_id is cEncoded UUID)
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
