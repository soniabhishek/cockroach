package auther

import (
	"github.com/crowdflux/angel/app/DAL/repositories/clients_repo"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/utilities"
	"github.com/gin-gonic/gin"
)

func GinAuther() gin.HandlerFunc {
	return func(c *gin.Context) {

		projectId, err := uuid.FromString(c.Param("projectId"))
		if utilities.IsValidError(err) {
			ShowAuthenticationErrorOverHttp(c, "project not valid")
			return
		}

		clientRepo := clients_repo.New()
		client, err := clientRepo.GetByProjectId(projectId)
		// Get Client ID
		if utilities.IsValidError(err) {
			ShowAuthenticationErrorOverHttp(c, "project not valid")
			return
		}

		if client.ClientSecretUuid == uuid.Nil {
			ShowAuthenticationErrorOverHttp(c, "project not valid")
			return
		}

		// Get Client Secret
		clientKey := c.Request.Header.Get("x-client-key")
		if clientKey == "" {
			ShowAuthenticationErrorOverHttp(c, "project not valid")
			return

		}

		if StdProdAuther.Check(client.ClientSecretUuid, clientKey) {

			//Let the request go forward, set a client_id param also
			c.Set("client_id", client.ID.String())
			c.Set("show_old", client.Options["show_old"])

		} else {
			ShowAuthenticationErrorOverHttp(c, "unauthorized")
		}

	}
}
