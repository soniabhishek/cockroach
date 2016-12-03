package work_flow_svc

import (
	"github.com/gin-gonic/gin"
	"work_flow_explorer_svc"
	"work_flow_io_svc"
	"work_flow_retriever_svc"
)

func AddHttpTransport(routerGroup *gin.RouterGroup) {
	work_flow_retriever_svc.AddHttpTransport(routerGroup)
	work_flow_io_svc.AddHttpTransport(routerGroup)
	work_flow_explorer_svc.AddHttpTransport(routerGroup)
}
