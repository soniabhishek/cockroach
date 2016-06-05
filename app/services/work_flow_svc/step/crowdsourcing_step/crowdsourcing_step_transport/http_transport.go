package crowdsourcing_step_transport

import "github.com/gin-gonic/gin"
import (
	"net/http"

	"gitlab.com/playment-main/angel/app/services/plerrors"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step/crowdsourcing_step"
)

func AddHttpTransport(r *gin.RouterGroup) {

	r.POST("/angel_crowdsourcing_gateway?action=flu_updates", crowdSourcingPostHandler())
}

type fluUpdateReq struct {
	FluUpdates []crowdsourcing_step.FluUpdate `json:"flu_updates"`
}

func crowdSourcingPostHandler() gin.HandlerFunc {

	return func(c *gin.Context) {

		var fluUpdateReq fluUpdateReq

		if err := c.Bind(&fluUpdateReq); err != nil {
			showErrorResponse(c, err)
			return
		}

		err := crowdsourcing_step.FluUpdateHandler(fluUpdateReq.FluUpdates)
		if err != nil {
			showErrorResponse(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

//--------------------------------------------------------------------------------//
//Helper

func showErrorResponse(c *gin.Context, err error) {

	var msg interface{}

	switch err.(type) {
	case plerrors.ServiceError:
		msg = err.(plerrors.ServiceError)
	case plerrors.IncorrectUUIDError:
		msg = err.(plerrors.IncorrectUUIDError)
	case plerrors.RequestParamMissingError:
		msg = err.(plerrors.RequestParamMissingError)

	//Commenting out the postgres driver error for now
	//case *pq.Error:
	//	msg = err.(*pq.Error)

	default:
		msg = err.Error()
	}
	c.JSON(http.StatusOK, gin.H{
		"error":   msg,
		"success": false,
	})
}
