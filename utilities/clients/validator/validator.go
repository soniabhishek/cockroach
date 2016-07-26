package validator

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/utilities"
	"gitlab.com/playment-main/angel/utilities/clients/models"
)

//--------------------------------------------------------------------------------//
//Validator

func ValidateInput(obj utilModels.Client) (err error) {

	// Validating UserName
	if obj.UserName == "" {
		err = errors.New("username invalid")
		return
	}

	// Validating Password
	if obj.Password == "" {
		err = errors.New("password invalid")
		return
	}

	// Validating ProjectLabel
	if obj.ProjectLabel == "" {
		err = errors.New("projectLabel invalid")
		return
	}

	// Validating ProjectName
	if obj.ProjectName == "" {
		err = errors.New("projectName invalid")
		return
	}

	// Validating URL
	if utilities.ValidateUrl(obj.Url) == false {
		err = errors.New("url invalid")
		return
	}

	// Validating Header
	if obj.Header == nil {
		err = errors.New("header invalid")
		return
	}

	return
}

//--------------------------------------------------------------------------------//
//Helpers

func ShowErrorResponse(err error) {
	plog.Info("Validation Error: ", err.Error())
}

func ShowErrorResponseOverHttp(c *gin.Context, err error) {

	var msg interface{}
	msg = err.Error()
	c.JSON(http.StatusOK, gin.H{
		"error":   msg,
		"success": false,
	})
}
