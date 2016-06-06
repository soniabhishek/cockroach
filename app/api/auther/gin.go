package auther

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/app/DAL/repositories/projects_repo"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/utilities"
)

func GinAuther() gin.HandlerFunc {
	return func(c *gin.Context) {

		projectId, err := uuid.FromString(c.Param("projectId"))

		projRepo := projects_repo.New()
		project, err := projRepo.GetById(projectId)
		// Get Client ID
		if utilities.IsValidError(err) {
			plog.Error("ClientId not found for ProjectId ["+projectId.String()+"]:", err)
			c.Header("authenication error", "client_id not found")
			c.AbortWithStatus(401)
			return
		}

		clId := project.ClientId.String()
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
