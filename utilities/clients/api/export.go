package utils_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/app/api/auther"
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/utilities/clients/models"
	"gitlab.com/playment-main/angel/utilities/clients/operations"
	"net/http"
)

const ENDPOINT = "/createuser"

var OK bool = true

func AddHttpTransport(routerGroup *gin.RouterGroup) {
	fmt.Println("Reached here!")
	routerGroup.POST(ENDPOINT, createUserHandler())
}

//--------------------------------------------------------------------------------//

func createUserHandler() gin.HandlerFunc {

	return func(c *gin.Context) {

		obj, err := validateInputFLU(c)
		if err != nil {
			// Incoming FLU is not valid.
			return
		}

		service := operations.Service{}

		_, err = service.CreateClient(&obj)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"userdetails": utilModels.Client{
					UserName:        obj.UserName,
					ClientId:        obj.ClientId,
					ClientSecretStr: auther.StdProdAuther.GetAPIKey(obj.ClientSecretId),
					ProjectId:       obj.ProjectId,
					ProjectLabel:    obj.ProjectLabel,
					ProjectName:     obj.ProjectName,
					Url:             obj.Url,
					Header:          obj.Header,
					Steps:           obj.Steps,

					Gender:    obj.Gender,
					FirstName: obj.FirstName,
					LastName:  obj.LastName,
					Phone:     obj.Phone,
				},
			})
		} else {
			plog.Error("Error while creating user: ", err)
			showErrorResponse(c, err)
		}

	}
}
