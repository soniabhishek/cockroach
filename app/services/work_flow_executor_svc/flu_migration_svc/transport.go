package flu_migration_svc

import "github.com/gin-gonic/gin"
import (
	"net/http"

	"errors"
	"fmt"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/services/plerrors"
	"github.com/lib/pq"
)

func AddHttpTransport(r *gin.RouterGroup) {

	fluMigrationService := New()

	r.POST("/flu_migration", migrationDetailsHandler(fluMigrationService))
}

func migrationDetailsHandler(fluMigrationService IFluMigrationService) gin.HandlerFunc {

	return func(c *gin.Context) {

		if c.Query("action") != "get_details" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": plerrors.ErrorBody{
					Code:    "FM_9999",
					Message: "Unsupported action. Only allowed action is 'get_details'",
				},
			})
			return
		}

		file, _, err := c.Request.FormFile("master_flu_csv")
		if err != nil {
			showErrorResponse(c, errors.New("problem getting uploaded file"))
			return
		}

		refName := c.PostForm("reference_name")
		if refName == "" {
			showErrorResponse(c, errors.New("provide reference_name"))
			return
		}

		fluMigrationCSVDetails, err := fluMigrationService.GetFluMigrationDetails(file, refName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": plerrors.ErrorBody{
					Code:    "FM_0001",
					Message: err.Error(),
					MetaData: models.JsonF{
						"raw_error": fmt.Sprint(err),
					},
				},
			})
			return
		}

		migrationDetails := models.JsonF{}
		msg := ""

		if fluMigrationCSVDetails.CrowdsourcingBufferDeleteFile != nil {
			migrationDetails["crowdsourcing_flu_buffer_delete_file"] = fluMigrationCSVDetails.CrowdsourcingBufferDeleteFileName
			msg += "Delete Flus from crowdsourcing flu buffer: https://api.playment.in/downloads/" + fluMigrationCSVDetails.CrowdsourcingBufferDeleteFile.Name() + " . "
		}

		if fluMigrationCSVDetails.UnificationBufferDeleteFile != nil {
			migrationDetails["unification_buffer_delete_file"] = fluMigrationCSVDetails.UnificationBufferDeleteFileName
			msg += "Delete Flus from unification flu buffer: https://api.playment.in/downloads/" + fluMigrationCSVDetails.UnificationBufferDeleteFile.Name() + " . "
		}

		if fluMigrationCSVDetails.DeactivateFluFile != nil {
			migrationDetails["deactivate_flu_file"] = fluMigrationCSVDetails.DeactivateFluFileName
			msg += "Flus to deactivate: https://api.playment.in/downloads/" + fluMigrationCSVDetails.DeactivateFluFile.Name() + " . "
		}

		if msg == "" {
			migrationDetails["message"] = "Nothing to clear."
		} else {
			migrationDetails["message"] = msg
		}

		c.JSON(http.StatusOK, gin.H{
			"success":           true,
			"migration_details": migrationDetails,
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
	case *pq.Error:
		msg = err.(*pq.Error)
	default:
		msg = err.Error()
	}
	c.JSON(http.StatusInternalServerError, plerrors.ErrorResponse{
		Success: false,
		Error: plerrors.ErrorBody{
			Code:    "FM_0002",
			Message: fmt.Sprint(msg),
		},
	})
}
