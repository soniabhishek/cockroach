package client_svc

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/support/app/services/flu_svc/flu_svc_transport"
)

func AddHttpTransport(r gin.RouterGroup) {
	flu_svc_transport.AddHttpTransport(r)
}
