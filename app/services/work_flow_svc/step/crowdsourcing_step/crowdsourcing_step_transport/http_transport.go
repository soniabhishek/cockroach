package crowdsourcing_step_transport

import "github.com/gin-gonic/gin"
import (
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/services/plerrors"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step/crowdsourcing_step"
	"net/http"
)

func AddHttpTransport(r *gin.RouterGroup) {

	r.POST("crowdsourcing/questionanswer", crowdSourcingPostHandler())
}

type markQuestionRequest struct {
	Question models.Question
	Answer   models.QuestionAnswer `json:"answer"`
}

func crowdSourcingPostHandler() gin.HandlerFunc {

	return func(c *gin.Context) {

		var markQuestionRequest markQuestionRequest

		if err := c.Bind(&markQuestionRequest); err != nil {
			showErrorResponse(c, err)
		}

		err := crowdsourcing_step.Std.HandleQuestionComplete(markQuestionRequest.Question, markQuestionRequest.Answer)
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
