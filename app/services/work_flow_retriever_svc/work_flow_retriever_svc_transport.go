package work_flow_retriever_svc

import (
	"github.com/gin-gonic/gin"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"net/http"
)

func AddHttpTransport(routerGroup *gin.RouterGroup) {
	workFlowRetrieverService := New()
	routerGroup.GET("/workflow", workFlowGetHandler(workFlowRetrieverService))

}

func workFlowGetHandler(workFlowService IWorkflowRetrieverService) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectId, err := uuid.FromString(c.Query("projectId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, "Invalid ProjectId")
			return
		}

		tag := c.Query("tag")
		if err != nil {
			c.JSON(http.StatusBadRequest, "Invalid tag")
			return
		}

		response, err := workFlowService.GetWorkFlow(projectId, tag)
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