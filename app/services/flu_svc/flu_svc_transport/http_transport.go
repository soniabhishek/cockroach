package flu_svc_transport

import (
	"net/http"
	"time"

	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services"
	"github.com/crowdflux/angel/app/services/flu_svc"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_errors"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_validator"
	"github.com/crowdflux/angel/app/services/plerrors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

//TODO Create another file for validator http transport. In future we may have to make a separate service for validatorss

func AddHttpTransport(routerGroup *gin.RouterGroup) {

	fluService := flu_svc.NewWithExposedValidators()

	routerGroup.POST("/project/:projectId/feedline", feedLineInputHandler(fluService))
	routerGroup.POST("project/:projectId/csv/feedline", csvFLUGenerator(fluService))

	routerGroup.GET("project/:projectId/upload/status", getUploadStatus(fluService))
	routerGroup.GET("/project/:projectId/feedline/:feedlineId", feedLineGetHandler(fluService))

	routerGroup.GET("/project/:projectId/validator", validatorGetHandler(fluService))
	routerGroup.POST("/project/:projectId/validator", validatorUpdateHandler(fluService))
}

//--------------------------------------------------------------------------------//

type fluPostResponse struct {
	Id          uuid.UUID `json:"flu_id"`
	ReferenceId string    `json:"reference_id"`
	Tag         string    `json:"tag"`
}

var requestLogger = plog.NewLogger("inbound_request", "INFO", "FILE")

//Inserts into mongo
func feedLineInputHandler(fluService flu_svc.IFluServiceExtended) gin.HandlerFunc {

	return func(c *gin.Context) {

		requestTime := time.Now()
		var flu models.FeedLineUnit

		var projectId uuid.UUID
		var err error
		projectId, err = uuid.FromString(c.Param("projectId"))
		var body []byte

		if err != nil {
			plog.Error("Invalid ProjectId from client", err, c.Param("projectId"))
			httpCode, resp := showErrorResponse(c, plerrors.ErrIncorrectUUID("projectId"))
			requestLogger.Info(formatLog(requestTime, body, time.Since(requestTime), httpCode, resp))
			return
		}
		// Validating JSON
		if err = c.BindJSON(&flu); err != nil {
			plog.Error("Error binding flu from client : ", err, "Body : ", body)
			httpCode, resp := showErrorResponse(c, plerrors.ErrMalformedJson)
			requestLogger.Info(formatLog(requestTime, body, time.Since(requestTime), httpCode, resp))
			return
		}
		flu.ProjectId = projectId

		plog.Info("Request Inbound. Reference: ", flu.ReferenceId, " Time : ", requestTime, " Tag : ", flu.Tag, " ProjectId : ", flu.ProjectId)

		err = fluService.AddFeedLineUnit(&flu)
		if err != nil {
			if err == projects_repo.ErrProjectNotFound {
				//Temporary hack. Wait for schema refactoring
				err = plerrors.ServiceError{"PR_0001", "Project not found"}
			}
			plog.Error("Error while adding flu to workflow ", err, flu)
			httpCode, resp := showErrorResponse(c, err)
			requestLogger.Info(formatLog(requestTime, body, time.Since(requestTime), httpCode, resp))
			return
		}

		// This has to be done for chutiya paytm dev
		if c.Keys["show_old"] == true {
			requestLogger.Info(formatLog(requestTime, body, time.Since(requestTime), http.StatusOK, models.JsonF{"success": true,
				"flu_id":       flu.ID,
				"reference_id": flu.ReferenceId,
				"tag":          flu.Tag}))

			c.JSON(http.StatusOK, gin.H{
				"success":      true,
				"flu_id":       flu.ID,
				"reference_id": flu.ReferenceId,
				"tag":          flu.Tag,
			})
		} else {
			fluResp := fluPostResponse{Id: flu.ID,
				ReferenceId: flu.ReferenceId,
				Tag:         flu.Tag}

			requestLogger.Info(formatLog(requestTime, body, time.Since(requestTime), http.StatusOK, models.JsonF{"success": true,
				"feed_line_unit": fluResp}))

			c.JSON(http.StatusOK, gin.H{
				"success":        true,
				"feed_line_unit": fluResp,
			})
		}
	}
}

func formatLog(requestTime time.Time, body []byte, duration time.Duration, httpCode int, response models.JsonF) []interface{} {

	return []interface{}{
		models.JsonF{"request_time": requestTime},
		models.JsonF{"body": body},
		models.JsonF{"request_duration": duration},
		models.JsonF{"http_code": httpCode},
		models.JsonF{"response": response},
	}

}

//This Handler is for Generating flus via csv upload
//csv should be 3 column wide as [reference_id, tag, build]
//This Handler will first Copy all content on the disk it can be varied upon on requirement to take directly in memory.
//Once file would be successfully written to disk a empty response would be sent back to user.
func csvFLUGenerator(fluService flu_svc.IFluServiceExtended) gin.HandlerFunc {
	return func(c *gin.Context) {

		//Validating ProjectId and checking if exist in database.
		projectId, err := uuid.FromString(c.Param("projectId"))
		if err != nil {
			plog.Error("Invalid ProjectId in CSV upload", err, c.Param("projectId"))
			services.SendBadRequest(c, "FLS000", err.Error(), nil)
			return
		}

		//Fetching descriptors for uploaded multipart file
		file, header, err := c.Request.FormFile("upload")
		if err != nil {
			plog.Error("Err", errors.New("problem in uploaded file"), err)
			services.SendBadRequest(c, "FLS002", err.Error(), nil)
			return
		}
		defer file.Close()

		filename := header.Filename

		plog.Info("Sent file for upload: ", filename)

		if err := fluService.CsvCheckBasicValidation(file, filename, projectId); err != nil {
			plog.Error("Err", errors.New("Validation Failed for csv"), err)
			services.SendBadRequest(c, "Fail", err.Error(), nil)
			return
		}

		//return response for file upload
		services.SendSuccessResponse(c, nil)
	}

}

func getUploadStatus(fluService flu_svc.IFluServiceExtended) gin.HandlerFunc {
	return func(c *gin.Context) {
		response, err := fluService.GetUploadStatus(c.Param("projectId"))
		if err != nil {
			services.SendSuccessResponse(c, "No Such ProjectID Exist")
		} else {
			services.SendSuccessResponse(c, response)
		}
	}
}

type fluGetResponse struct {
	Id          uuid.UUID    `json:"id"`
	ReferenceId string       `json:"reference_id"`
	Tag         string       `json:"tag"`
	Data        models.JsonF `json:"data"`
	CreatedAt   time.Time    `json:"created_at"`
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
			showErrorResponse(c, flu_errors.ErrTagMissing)
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

func showErrorResponse(c *gin.Context, err error) (statusCode int, resp models.JsonF) {

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

	statusCode = http.StatusOK
	if err == flu_errors.ErrRequestTimedOut {
		statusCode = http.StatusGatewayTimeout
	}
	resp = models.JsonF{"error": msg,
		"success": false}

	c.JSON(statusCode, gin.H{
		"error":   msg,
		"success": false,
	})
	return
}
