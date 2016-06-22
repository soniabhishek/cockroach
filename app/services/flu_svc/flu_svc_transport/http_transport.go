package flu_svc_transport

import (
	"net/http"
	"time"

	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/app/DAL/repositories/projects_repo"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/services/flu_svc"
	"gitlab.com/playment-main/angel/app/services/flu_svc/flu_validator"
	"gitlab.com/playment-main/angel/app/services/plerrors"
)

//TODO Create another file for validator http transport. In future we may have to make a separate service for validatorss

func AddHttpTransport(routerGroup *gin.RouterGroup) {

	fluService := flu_svc.NewWithExposedValidators()

	routerGroup.POST("/project/:projectId/feedline", feedLineInputHandler(fluService))
	routerGroup.GET("/project/:projectId/feedline/:feedlineId", feedLineGetHandler(fluService))

	routerGroup.GET("/project/:projectId/validator", validatorGetHandler(fluService))
	routerGroup.POST("/project/:projectId/validator", validatorUpdateHandler(fluService))
}

//--------------------------------------------------------------------------------//

type fluPostResponse struct {
	Id          uuid.UUID `json:"id"`
	ReferenceId string    `json:"reference_id"`
	Tag         string    `json:"tag"`
}

func feedLineInputHandler(fluService flu_svc.IFluService) gin.HandlerFunc {

	return func(c *gin.Context) {

		flu, err := validateInputFLU(c, fluService)
		if err != nil {
			// Incoming FLU is not valid.
			return
		}

		// This has to be done for chutiya paytm dev
		if c.Keys["show_old"] == true {
			c.JSON(http.StatusOK, gin.H{
				"success":      true,
				"flu_id":       flu.ID,
				"reference_id": flu.ReferenceId,
				"tag":          flu.Tag,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"feed_line_unit": fluPostResponse{
					Id:          flu.ID,
					ReferenceId: flu.ReferenceId,
					Tag:         flu.Tag,
				},
			})
		}
	}
}

type fluGetResponse struct {
	Id          uuid.UUID       `json:"id"`
	ReferenceId string          `json:"reference_id"`
	Tag         string          `json:"tag"`
	Data        models.JsonFake `json:"data"`
	CreatedAt   time.Time       `json:"created_at"`
}

func feedLineGetHandler(fluService flu_svc.IFluService) gin.HandlerFunc {

	return func(c *gin.Context) {
		feedLineId, err := uuid.FromString(c.Param("feedlineId"))
		if err != nil {
			showErrorResponse(c, plerrors.ErrIncorrectUUID("feedlineId"))
			return
		}
		flu, err := fluService.GetFeedLineUnit(feedLineId)
		if err != nil {
			showErrorResponse(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"feed_line_unit": fluGetResponse{
				Id:          flu.ID,
				ReferenceId: flu.ReferenceId,
				Tag:         flu.Tag,
				Data:        flu.Data,
				CreatedAt:   flu.CreatedAt.Time,
			},
			"success": true,
		})
	}
}

//--------------------------------------------------------------------------------//]
//Validator Service part
//Need to move to another file

type validatorGetResponse struct {
	FieldName   string `json:"field_name"`
	Type        string `json:"type"`
	IsMandatory bool   `json:"is_mandatory"`
}

func validatorGetHandler(validatorSvc flu_validator.IFluValidatorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectId, err := uuid.FromString(c.Param("projectId"))
		if err != nil {
			showErrorResponse(c, plerrors.ErrIncorrectUUID("projectId"))
			return
		}
		tag := c.Query("tag")

		if tag == "" {
			showErrorResponse(c, flu_svc.ErrTagMissing)
			return
		}

		fvs, err := validatorSvc.GetValidators(projectId, tag)
		if err != nil {
			showErrorResponse(c, err)
			return
		}

		var fluValidatorList []validatorGetResponse

		for _, fv := range fvs {
			fluValidatorList = append(fluValidatorList, validatorGetResponse{
				FieldName:   fv.FieldName,
				Type:        fv.Type,
				IsMandatory: fv.IsMandatory,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"success":    true,
			"validators": fluValidatorList,
		})
	}
}

func validatorUpdateHandler(validatorSvc flu_validator.IFluValidatorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectId, err := uuid.FromString(c.Param("projectId"))
		if err != nil {
			showErrorResponse(c, plerrors.ErrIncorrectUUID("projectId"))
			return
		}

		var fv models.FLUValidator

		if err := c.BindJSON(&fv); err != nil {
			showErrorResponse(c, plerrors.ErrMalformedJson)
			return
		}

		fv.ProjectId = projectId
		err = validatorSvc.SaveValidator(&fv)
		if err != nil {
			showErrorResponse(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"validator": validatorGetResponse{
				FieldName:   fv.FieldName,
				IsMandatory: fv.IsMandatory,
				Type:        fv.Type,
			},
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
	case flu_validator.DataValidationError:
		msg = err.(flu_validator.DataValidationError)
	case plerrors.IncorrectUUIDError:
		msg = err.(plerrors.IncorrectUUIDError)
	case plerrors.RequestParamMissingError:
		msg = err.(plerrors.RequestParamMissingError)

	//Commenting out the postgres driver error for now
	//case *pq.Error:
	//	msg = err.(*pq.Error)

	default:
		msg = err.Error()
	}
	c.JSON(http.StatusOK, gin.H{
		"error":   msg,
		"success": false,
	})
}

//--------------------------------------------------------------------------------//
//Validator

func validateInputFLU(c *gin.Context, fluService flu_svc.IFluService) (flu models.FeedLineUnit, err error) {

	//Variable name will be changed to projectId after the schema refactoring
	var projectId uuid.UUID
	projectId, err = uuid.FromString(c.Param("projectId"))
	if err != nil {
		showErrorResponse(c, plerrors.ErrIncorrectUUID("projectId"))
		return
	}

	// Validating JSON
	if err = c.BindJSON(&flu); err != nil {
		showErrorResponse(c, plerrors.ErrMalformedJson)
		return
	}

	// Validating ReferenceID
	if flu.ReferenceId == "" {
		showErrorResponse(c, flu_svc.ErrReferenceIdMissing)
		return
	}

	// Validating TAG
	if flu.Tag == "" {
		err = flu_svc.ErrTagMissing
		showErrorResponse(c, flu_svc.ErrTagMissing)
		return
	}

	// Validating Data
	if flu.Data == nil {
		err = flu_svc.ErrDataMissing
		showErrorResponse(c, flu_svc.ErrDataMissing)
		return
	}

	flu.ProjectId = projectId
	err = fluService.AddFeedLineUnit(&flu)
	fmt.Println("Error while adding flu: ", err)
	if err != nil {
		if err == projects_repo.ErrProjectNotFound {
			//Temporary hack. Wait for schema refactoring
			err = plerrors.ServiceError{"PR_0001", "Project not found"}
		}
		showErrorResponse(c, err)
		return
	}
	return
}
