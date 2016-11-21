package user_svc

import (
	"net/http"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/utilities/clients/validator"
	"github.com/gin-gonic/gin"
)

func AddHttpTransport(routerGroup *gin.RouterGroup) {
	userService := New()
	routerGroup.POST("/users/create", createUserHandler(userService))
}

//--------------------------------------------------------------------------------//

func createUserHandler(userService IUserService) gin.HandlerFunc {

	return func(c *gin.Context) {

		obj := models.User{}
		err := c.BindJSON(&obj)
		if err != nil {
			validator.ShowErrorResponseOverHttp(c, err)
			return
		}
		err = validator.ValidateUser(obj)
		if err != nil {
			validator.ShowErrorResponseOverHttp(c, err)
			return
		}

		err = userService.CreateUser(obj)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
			})
		} else {
			plog.Error("Error while creating user: ", err)
			validator.ShowErrorResponseOverHttp(c, err)
		}

	}
}
