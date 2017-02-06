package project_svc

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/plog/log_tags"
	"github.com/crowdflux/angel/app/services"
	"github.com/gin-gonic/gin"
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
			plog.Error("Project_svc", err, plog.Message("fetchProjectsHandler. Invalid ClientId"), plog.MessageWithParam(log_tags.CLIENT_ID, c.Param("clientId")))

			services.SendBadRequest(c, "FETCHPROJECT", err.Error(), nil)
			return
		}
		response, err := workFlowProjectService.FetchProjectsByClientId(clientId)
		if err != nil {
			plog.Error("Project_svc", err, plog.Message("fetchProjectsHandler. Fetching Client Projects Error. ClientId "), plog.MessageWithParam(log_tags.CLIENT_ID, clientId))
			services.SendFailureResponse(c, "FETCHPROJECT", err.Error(), nil)
			return
		}
		services.SendSuccessResponse(c, response)
	}
}

func createProjectsHandler(workFlowProjectService IWorkFlowProjetService) gin.HandlerFunc {
	return func(c *gin.Context) {
		creatorId, ok := c.Get("userId")
		if !ok {
			services.SendBadRequest(c, "CREATEPROJECT", "NO Creator Id Present", nil)
			return
		}
		obj := models.Project{}
		err := c.BindJSON(&obj)
		if err != nil {
			services.SendBadRequest(c, "CREATEPROJECT", err.Error(), nil)
			return
		}

		obj.CreatorId, err = uuid.FromString(creatorId.(string))
		if err != nil {
			services.SendBadRequest(c, "CREATEPROJECT", err.Error(), nil)
			return
		}

		response, err := workFlowProjectService.CreateProject(obj)
		if err != nil {
			services.SendFailureResponse(c, "FETCHPROJECT", err.Error(), nil)
			plog.Error("Project_svc", err, plog.Message("CreateProjectHandler. Creating client Error. CreatorId"), plog.MessageWithParam(log_tags.USER_ID, obj.CreatorId))
			return
		}
		services.SendSuccessResponse(c, response)
	}
}
