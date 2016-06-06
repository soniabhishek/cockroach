package auther

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/app/DAL/repositories/clients_repo"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/utilities"
)

func GinAuther() gin.HandlerFunc {
	return func(c *gin.Context) {

		projectId, err := uuid.FromString(c.Param("projectId"))
		if utilities.IsValidError(err) {
			plog.Info("Auther", "projectId not uuid")
			c.Header("authenication error", "project not valid")
			c.AbortWithStatus(401)
			return
		}

		clientRepo := clients_repo.New()
		client, err := clientRepo.GetByProjectId(projectId)
		// Get Client ID
		if utilities.IsValidError(err) {
			plog.Error("ClientId not found for ProjectId ["+projectId.String()+"]:", err)
			c.Header("authenication error", "project not valid")
			c.AbortWithStatus(401)
			return
		}

		if client.ClientSecretUuid == uuid.Nil {
			plog.Error("ClientSecretID not found ["+client.ID.String()+"]:", err)
			c.Header("authenication error", "project not valid")
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

		if StdProdAuther.Check(client.ClientSecretUuid, clSec) {

			//Let the request go forward, set a user_id param also
			c.Set("client_id", client.ID.String())
		} else {
			c.Header("authenication error", "unauthorized")
			c.AbortWithStatus(401)
		}

	}
}
