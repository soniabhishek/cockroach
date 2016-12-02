package work_flow_explorer_svc

import (
	"net/http"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/utilities/clients/validator"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func AddHttpTransport(routerGroup *gin.RouterGroup) {
	workFlowExplorerService := New()
	routerGroup.GET("/clients", fetchClientsHandler(workFlowExplorerService))
	routerGroup.GET("/clients/:clientId/projects", fetchProjectsHandler(workFlowExplorerService))
	routerGroup.GET("/projects/:projectId/workflows", fetchWorkflowsHandler(workFlowExplorerService))
	routerGroup.POST("/clients", createClientHandler(workFlowExplorerService))
	routerGroup.POST("/clients/:clientId/projects", createProjectsHandler(workFlowExplorerService))
}

//--------------------------------------------------------------------------------//

func createClientHandler(workFlowExplorerService IWorkFlowExplorerService) gin.HandlerFunc {

	return func(c *gin.Context) {

		obj := models.Client{}
		err := c.BindJSON(&obj)
		if err != nil {
			validator.ShowErrorResponseOverHttp(c, err)
			return
		}
		client, err := workFlowExplorerService.CreateClient(obj)
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

func fetchClientsHandler(workFlowExplorerService IWorkFlowExplorerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		response, err := workFlowExplorerService.FetchAllClient()
		if err != nil {
			plog.Error("Fetching Clients Error", err)
			validator.ShowErrorResponseOverHttp(c, err)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"clients": response,
			})
		}

	}
}

func fetchProjectsHandler(workFlowExplorerService IWorkFlowExplorerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientId, err := uuid.FromString(c.Param("clientId"))
		if err != nil {
			plog.Error("Invalid Id", err)
			validator.ShowErrorResponseOverHttp(c, err)
			return
		}
		response, err := workFlowExplorerService.FetchProjectsByClientId(clientId)
		if err != nil {
			plog.Error("Fetching Client Projects Error", err)
			validator.ShowErrorResponseOverHttp(c, err)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success":  true,
				"projects": response,
			})
		}
	}
}

func fetchWorkflowsHandler(workFlowExplorerService IWorkFlowExplorerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectId, err := uuid.FromString(c.Param("projectId"))
		if err != nil {
			plog.Error("Invalid Id", err)
			validator.ShowErrorResponseOverHttp(c, err)
			return
		}
		response, err := workFlowExplorerService.FetchWorkflowsByProjectId(projectId)
		if err != nil {
			plog.Error("Fetching Projects workflows Error", err)
			validator.ShowErrorResponseOverHttp(c, err)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success":   true,
				"workflows": response,
			})
		}
	}
}

func createProjectsHandler(workFlowExplorerService IWorkFlowExplorerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		creatorId, ok := c.Get("userId")
		if !ok {
			err := errors.New("NO Creator Id Present")
			validator.ShowErrorResponseOverHttp(c, err)
			return
		}
		obj := models.Project{}
		err := c.BindJSON(&obj)
		if err != nil {
			validator.ShowErrorResponseOverHttp(c, err)
			return
		}

		obj.CreatorId, err = uuid.FromString(creatorId.(string))
		if err != nil {
			validator.ShowErrorResponseOverHttp(c, err)
			return
		}

		client, err := workFlowExplorerService.CreateProject(obj)
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
