package client_svc

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services"
	"github.com/gin-gonic/gin"
)

func AddHttpTransport(routerGroup *gin.RouterGroup) {
	workFlowClientService := New()
	routerGroup.GET("/clients", fetchClientsHandler(workFlowClientService))
	routerGroup.POST("/clients", createClientHandler(workFlowClientService))
}

func createClientHandler(workFlowClientService IWorkFlowClientService) gin.HandlerFunc {

	return func(c *gin.Context) {

		obj := models.Client{}
		err := c.BindJSON(&obj)
		if err != nil {
			services.SendBadRequest(c, "CREATECLIENT", err.Error(), nil)
			return
		}
		client, err := workFlowClientService.CreateClient(obj)
		if err != nil {
			services.SendFailureResponse(c, "CREATECLIENT", err.Error(), nil)
			plog.Error("Creating client Error", err)
			return
		}
		services.SendSuccessResponse(c, client)
	}
}

func fetchClientsHandler(workFlowClientService IWorkFlowClientService) gin.HandlerFunc {
	return func(c *gin.Context) {
		response, err := workFlowClientService.FetchAllClient()
		if err != nil {
			plog.Error("Fetching Clients Error", err)
			services.SendFailureResponse(c, "FETCHCLIENT", err.Error(), nil)
			return
		}
		services.SendSuccessResponse(c, response)
	}
}
