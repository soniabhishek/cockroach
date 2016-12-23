package services

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/services/plerrors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func SendSuccessResponse(c *gin.Context, response interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"success": true,
	})
}

func SendBadRequest(c *gin.Context, code, message string, data models.JsonF) {
	c.JSON(http.StatusBadRequest, plerrors.ErrorResponse{
		Success: false,
		Error: plerrors.ErrorBody{
			Code:     code,
			Message:  message,
			MetaData: data,
		},
	})
}

func SendFailureResponse(c *gin.Context, code, message string, data models.JsonF) {
	c.JSON(http.StatusExpectationFailed, plerrors.ErrorResponse{
		Success: false,
		Error: plerrors.ErrorBody{
			Code:     code,
			Message:  message,
			MetaData: data,
		},
	})
}

func AtoiOrDefault(s string, defaultVal int) (integer int) {
	integer, err := strconv.Atoi(s)
	if err != nil {
		integer = defaultVal
	}
	return
}
