package work_flow_svc

import (
	"client_svc"
	"github.com/gin-gonic/gin"
	"project_svc"
	"work_flow_io_svc"
)

func AddHttpTransport(routerGroup *gin.RouterGroup) {
	work_flow_io_svc.AddHttpTransport(routerGroup)
	client_svc.AddHttpTransport(routerGroup)
	project_svc.AddHttpTransport(routerGroup)
}
