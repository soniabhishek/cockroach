package utils_api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/utilities/clients/models"
	"net/http"
)

//--------------------------------------------------------------------------------//
//Helper

func showErrorResponse(c *gin.Context, err error) {

	var msg interface{}
	msg = err.Error()
	c.JSON(http.StatusOK, gin.H{
		"error":   msg,
		"success": false,
	})
}

//--------------------------------------------------------------------------------//
//Validator

func validateInputFLU(c *gin.Context) (obj utilModels.Client, err error) {

	// Validating JSON
	if err = c.BindJSON(&obj); err != nil {
		err = errors.New("Invalid JSON : " + err.Error())
		showErrorResponse(c, err)
		return
	}

	// Validating UserName
	if obj.UserName == "" {
		err = errors.New("username invalid")
		showErrorResponse(c, err)
		return
	}

	// Validating Password
	if obj.Password == "" {
		err = errors.New("password invalid")
		showErrorResponse(c, err)
		return
	}

	// Validating ProjectLabel
	if obj.ProjectLabel == "" {
		err = errors.New("projectLabel invalid")
		showErrorResponse(c, err)
		return
	}

	// Validating ProjectName
	if obj.ProjectName == "" {
		err = errors.New("projectName invalid")
		showErrorResponse(c, err)
		return
	}

	// Validating URL
	if obj.Url == "" {
		err = errors.New("url invalid")
		showErrorResponse(c, err)
		return
	}

	// Validating Header
	if obj.Header == nil {
		err = errors.New("header invalid")
		showErrorResponse(c, err)
		return
	}

	return
}
