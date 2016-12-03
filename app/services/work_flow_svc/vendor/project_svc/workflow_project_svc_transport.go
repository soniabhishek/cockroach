package project_svc

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
	workFlowProjectService := New()
	routerGroup.GET("/clients/:clientId/projects", fetchProjectsHandler(workFlowProjectService))
	routerGroup.POST("/clients/:clientId/projects", createProjectsHandler(workFlowProjectService))
}

func fetchProjectsHandler(workFlowProjectService IWorkFlowProjetService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientId, err := uuid.FromString(c.Param("clientId"))
		if err != nil {
			plog.Error("Invalid Id", err)
			validator.ShowErrorResponseOverHttp(c, err)
			return
		}
		response, err := workFlowProjectService.FetchProjectsByClientId(clientId)
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

func createProjectsHandler(workFlowProjectService IWorkFlowProjetService) gin.HandlerFunc {
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

		client, err := workFlowProjectService.CreateProject(obj)
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
