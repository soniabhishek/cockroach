package work_flow_io_svc

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddHttpTransport(routerGroup *gin.RouterGroup) {
	workFlowBuilderService := New()

	routerGroup.GET("/workflow/:workflowId", workFlowGetHandler(workFlowBuilderService))
	routerGroup.POST("/workflow", addWorkFlowHandler(workFlowBuilderService))
	routerGroup.POST("/workflow/init/:projectId", intitWorkFlowHandler(workFlowBuilderService))
	routerGroup.PUT("/workflow", updateWorkFlowHandler(workFlowBuilderService))

}

func workFlowGetHandler(workFlowService IWorkflowBuilderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowId, err := uuid.FromString(c.Param("workflowId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, "Invalid id")
			return
		}

		response, err := workFlowService.GetWorkflowContainer(workflowId)
		if err != nil {
			plog.Error("WorkflowFetching : ", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Failure",
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data":    response,
			"success": true,
		})
	}

}

func intitWorkFlowHandler(workFlowService IWorkflowBuilderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID, err := uuid.FromString(c.Param("projectId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, "Invalid id")
			return
		}

		response, err := workFlowService.InitWorkflowContainer(projectID)
		if err != nil {
			plog.Error("WorkflowInitializing : ", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Failure",
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, response)
	}
}

func addWorkFlowHandler(workFlowService IWorkflowBuilderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowContainer := models.WorkflowContainer{}
		if err := c.BindJSON(&workflowContainer); err != nil {
			c.JSON(http.StatusBadRequest, "Invalid workflow "+err.Error())
			return
		}
		response, err := workFlowService.AddWorkflowContainer(workflowContainer)
		if err != nil {
			plog.Error("WorkflowAdd : ", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Failure",
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, response)
	}
}

func updateWorkFlowHandler(workFlowService IWorkflowBuilderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowContainer := models.WorkflowContainer{}
		if err := c.BindJSON(&workflowContainer); err != nil {
			c.JSON(http.StatusBadRequest, "Invalid workflow"+err.Error())
			return
		}

		response, err := workFlowService.UpdateWorkflowContainer(workflowContainer)
		if err != nil {
			plog.Error("WorkflowUpdate : ", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Failure",
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data":    response,
			"success": true,
		})
	}
}
