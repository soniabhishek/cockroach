package work_flow_io_svc

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services"
	"github.com/gin-gonic/gin"
)

func AddHttpTransport(routerGroup *gin.RouterGroup) {
	workFlowBuilderService := New()

	routerGroup.GET("/workflows", fetchWorkflowsHandler(workFlowBuilderService))
	routerGroup.GET("/workflows/:workflowId", workFlowGetHandler(workFlowBuilderService))
	routerGroup.PUT("/workflows", updateWorkFlowHandler(workFlowBuilderService))
	routerGroup.POST("/workflows", addWorkFlowHandler(workFlowBuilderService))
	routerGroup.POST("/workflows/:workflowId/action/clone", workFlowClonehandler(workFlowBuilderService))
}

/**
This is used to return the Workflow based on the provided workflowId
*/
func workFlowGetHandler(workFlowService IWorkflowBuilderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowId, err := uuid.FromString(c.Param("workflowId"))
		if err != nil {
			services.SendBadRequest(c, "WFGET", "Invalid ID", nil)
			return
		}
		response, err := workFlowService.GetWorkflowContainer(workflowId)
		if err != nil {
			plog.Error("Work_flow_io_svc", err, plog.NewMessageWithParam("workFlowGetHandler. Error fetching workflow. work_flow_id", workflowId))
			services.SendFailureResponse(c, "WFGET", err.Error(), nil)
			return
		}
		services.SendSuccessResponse(c, response)
	}
}

func addWorkFlowHandler(workFlowService IWorkflowBuilderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowContainer := models.WorkflowContainer{}
		if err := c.BindJSON(&workflowContainer); err != nil {
			services.SendBadRequest(c, "WFADD", "Invalid workflow "+err.Error(), nil)
			return
		}
		response, err := workFlowService.AddWorkflowContainer(workflowContainer)
		if err != nil {
			plog.Error("Work_flow_io_svc", err, plog.NewMessage("addWorkFlowHandler. Error adding workflow"))
			services.SendFailureResponse(c, "WFADD", err.Error(), nil)
			return
		}
		services.SendSuccessResponse(c, response)
	}
}

func updateWorkFlowHandler(workFlowService IWorkflowBuilderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowContainer := models.WorkflowContainer{}
		if err := c.BindJSON(&workflowContainer); err != nil {
			services.SendBadRequest(c, "WFUPDATE", "Invalid workflow "+err.Error(), nil)
			return
		}
		response, err := workFlowService.UpdateWorkflowContainer(workflowContainer)
		if err != nil {
			plog.Error("Work_flow_io_svc", err, plog.NewMessage("updateWorkFlowHandler. Error updating workflow"))
			services.SendFailureResponse(c, "WFUPDATE", err.Error(), nil)
			return
		}
		services.SendSuccessResponse(c, response)
	}
}

func workFlowClonehandler(workFlowService IWorkflowBuilderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowId, err := uuid.FromString(c.Param("workflowId"))
		if err != nil {
			services.SendBadRequest(c, "WFCLONE", "Invalid Id", nil)
			return
		}
		workflowCloneData := models.WorkFlowCloneModel{}
		if err := c.BindJSON(&workflowCloneData); err != nil {
			services.SendBadRequest(c, "WFCLONE", "Invalid Clone Request "+err.Error(), nil)
			return
		}
		workflowCloneData.WorkFlowId = workflowId
		response, err := workFlowService.CloneWorkflowContainer(workflowCloneData)
		if err != nil {
			services.SendFailureResponse(c, "WFCLONE", err.Error(), nil)
			return
		}
		services.SendSuccessResponse(c, response)
	}
}

func fetchWorkflowsHandler(workFlowService IWorkflowBuilderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectId, err := uuid.FromString(c.Query("projectId"))
		if err != nil {
			services.SendBadRequest(c, "WFFETCH", "Invalid ProjectId", nil)
			return
		}

		tag := c.Query("tag")

		response, err := workFlowService.FetchWorkflows(projectId, tag)
		if err != nil {
			plog.Error("Work_flow_io_svc", err, plog.NewMessage("fetchWorkflowsHandler. Fetching Projects workflows Error"), plog.NewMessageWithParam("tag", tag))
			services.SendFailureResponse(c, "WFFETCH", err.Error(), nil)
			return
		}
		services.SendSuccessResponse(c, response)
	}
}
