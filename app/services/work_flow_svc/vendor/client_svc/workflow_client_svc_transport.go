package client_svc

import (
	"net/http"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/utilities/clients/validator"
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
			validator.ShowErrorResponseOverHttp(c, err)
			return
		}
		client, err := workFlowClientService.CreateClient(obj)
		if err != nil {
			validator.ShowErrorResponseOverHttp(c, err)
			plog.Error("Creating client Error", err)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    client,
			})
		}

	}
}

func fetchClientsHandler(workFlowClientService IWorkFlowClientService) gin.HandlerFunc {
	return func(c *gin.Context) {
		response, err := workFlowClientService.FetchAllClient()
		if err != nil {
			plog.Error("Fetching Clients Error", err)
			validator.ShowErrorResponseOverHttp(c, err)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    response,
			})
		}

	}
}
