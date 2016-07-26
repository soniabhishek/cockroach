package utils_api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/utilities/clients/models"
)

//--------------------------------------------------------------------------------//
//Validator
func validateJson(c *gin.Context) (obj utilModels.Client, err error) {

	// Validating JSON
	if err = c.BindJSON(&obj); err != nil {
		err = errors.New("Invalid JSON : " + err.Error())
		return
	}
	return
}
