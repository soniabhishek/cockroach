package work_flow_explorer_svc

import (
	"net/http"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/utilities/clients/validator"
	"github.com/gin-gonic/gin"
)

const ENDPOINT = "/create-client"
const FETCHCLIENTS = "/fetch-clients"
const FETCHPROJECTS = "/fetch-client-projects/:clientId"

func AddHttpTransport(routerGroup *gin.RouterGroup) {
	clientService := New()
	routerGroup.GET(FETCHCLIENTS, fetchClientsHandler(clientService))
	routerGroup.GET(FETCHPROJECTS, fetchProjectsHandler(clientService))
	routerGroup.POST(ENDPOINT, createClientHandler(clientService))
}

//--------------------------------------------------------------------------------//

func createClientHandler(clientService IWorkFlowExplorerService) gin.HandlerFunc {

	return func(c *gin.Context) {

		obj := models.Client{}
		err := c.BindJSON(&obj)
		if err != nil {
			validator.ShowErrorResponseOverHttp(c, err)
			return
		}
		err = validator.ValidateClient(obj)
		if err != nil {
			validator.ShowErrorResponseOverHttp(c, err)
			return
		}
		client, err := clientService.CreateClient(obj)
		if err != nil {
			validator.ShowErrorResponseOverHttp(c, err)
			plog.Error("Creating client Error", err)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"client":  client,
			})
		}

	}
}

func fetchClientsHandler(clientService IWorkFlowExplorerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		response, err := clientService.FetchAllClient()
		if err != nil {
			plog.Error("Fetching Clients Error", err)
			validator.ShowErrorResponseOverHttp(c, err)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"clients": response,
			})
		}

	}
}

func fetchProjectsHandler(clientService IWorkFlowExplorerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientId, err := uuid.FromString(c.Param("clientId"))
		if err != nil {
			plog.Error("Invalid Id", err)
			validator.ShowErrorResponseOverHttp(c, err)
		}
		response, err := clientService.FetchProjectsByClientId(clientId)
		if err != nil {
			plog.Error("Fetching Client Projects Error", err)
			validator.ShowErrorResponseOverHttp(c, err)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success":  true,
				"projects": response,
			})
		}
	}
}
