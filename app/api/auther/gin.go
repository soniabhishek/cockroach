package auther

import (
	"github.com/crowdflux/angel/app/DAL/repositories/clients_repo"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/utilities"
	"github.com/gin-gonic/gin"
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
		clientKey := c.Request.Header.Get("x-client-key")
		if clientKey == "" {
			c.Header("authenication error", "x-client-key required")
			c.AbortWithStatus(401)
			return

		}

		if StdProdAuther.Check(client.ClientSecretUuid, clientKey) {

			//Let the request go forward, set a client_id param also
			c.Set("client_id", client.ID.String())
			c.Set("show_old", client.Options["show_old"])

		} else {
			c.Header("authenication error", "unauthorized")
			c.AbortWithStatus(401)
		}

	}
}
