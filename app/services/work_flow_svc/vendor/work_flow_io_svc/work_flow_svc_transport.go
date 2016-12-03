package work_flow_io_svc

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/utilities/clients/validator"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddHttpTransport(routerGroup *gin.RouterGroup) {
	workFlowBuilderService := New()

	routerGroup.GET("/workflows", fetchWorkflowsHandler(workFlowBuilderService))
	routerGroup.GET("/workflows/:workflowId", workFlowGetHandler(workFlowBuilderService))
	routerGroup.POST("/workflows", addWorkFlowHandler(workFlowBuilderService))
	routerGroup.PUT("/workflows", updateWorkFlowHandler(workFlowBuilderService))
	routerGroup.POST("/workflows/:workflowId/action/clone", workFlowClonehandler(workFlowBuilderService))
}

/**
This is used to return the Workflow based on the provided workflowId
*/
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
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data":    response,
			"success": true,
		})
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
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data":    response,
			"success": true,
		})
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
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data":    response,
			"success": true,
		})
	}
}

func workFlowClonehandler(workFlowService IWorkflowBuilderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		workflowId, err := uuid.FromString(c.Param("workflowId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, "Invalid id")
			return
		}
		workflowCloneData := models.WorkFlowCloneModel{}
		if err := c.BindJSON(&workflowCloneData); err != nil {
			c.JSON(http.StatusBadRequest, "Invalid Clone Request "+err.Error())
			return
		}
		workflowCloneData.WorkFlowId = workflowId
		response, err := workFlowService.CloneWorkflowContainer(workflowCloneData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data":    response,
			"success": true,
		})
	}
}

func fetchWorkflowsHandler(workFlowService IWorkflowBuilderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectId, err := uuid.FromString(c.Param("projectId"))
		if err != nil {
			plog.Error("Invalid Id", err)
			validator.ShowErrorResponseOverHttp(c, err)
			return
		}
		response, err := workFlowService.FetchWorkflowsByProjectId(projectId)
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
