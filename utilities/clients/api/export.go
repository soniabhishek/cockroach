package utils_api

import (
	"fmt"
	"net/http"

	"github.com/crowdflux/angel/app/api/auther"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/utilities/clients/models"
	"github.com/crowdflux/angel/utilities/clients/operations"
	"github.com/crowdflux/angel/utilities/clients/validator"
	"github.com/gin-gonic/gin"
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

		obj, err := validateJson(c)
		if err != nil {
			validator.ShowErrorResponseOverHttp(c, err)
			return
		}
		err = validator.ValidateInput(obj)
		if err != nil {
			validator.ShowErrorResponseOverHttp(c, err)
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
			validator.ShowErrorResponseOverHttp(c, err)
		}

	}
}
