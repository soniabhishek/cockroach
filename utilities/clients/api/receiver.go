package utils_api

import (
	"errors"

	"github.com/crowdflux/angel/utilities/clients/models"
	"github.com/gin-gonic/gin"
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
