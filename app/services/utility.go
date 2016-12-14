package services

import (
	"github.com/crowdflux/angel/app/services/plerrors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendSuccessResponse(c *gin.Context, response interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"success": true,
	})
}

func SendBadRequest(c *gin.Context, code, message string, data interface{}) {
	c.JSON(http.StatusBadRequest, plerrors.ErrorResponse{
		Success: false,
		Error: plerrors.ErrorBody{
			Code:     code,
			Message:  message,
			MetaData: data,
		},
	})
}

func SendFailureResponse(c *gin.Context, code, message string, data interface{}) {
	c.JSON(http.StatusExpectationFailed, plerrors.ErrorResponse{
		Success: false,
		Error: plerrors.ErrorBody{
			Code:     code,
			Message:  message,
			MetaData: data,
		},
	})
}
