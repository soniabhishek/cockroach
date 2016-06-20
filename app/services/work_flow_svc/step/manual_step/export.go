package manual_step

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
	"io"
	"net/http"
	"os"
)

var StdManualStep = &manualStep{
	Step: step.New(),
}

// Just a short form for above
var Std = StdManualStep

func AddHttpTransport(r *gin.RouterGroup) {

	r.POST("/manual_step_download", fileDownloadHandler())
	r.POST("/manual_step_upload", fileUploadHandler())
}

func fileDownloadHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		plog.Info(uuid.NewV4().String())
		manualStepId, err := uuid.FromString(c.Request.Header.Get("manualStepId"))
		plog.Info(c.Param("manualStepId"), manualStepId, err)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "please provide params",
			})
			return
		}

		filePath := DownloadCsv(manualStepId)

		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"filepath": "/downloadedfiles" + string(os.PathSeparator) + manualStepId.String() + ".csv",
		})
		plog.Info(c.Request.RequestURI, filePath)
	}
}
func fileUploadHandler() gin.HandlerFunc {

	return func(c *gin.Context) {

		file, header, err := c.Request.FormFile("upload")
		if err != nil {
			plog.Error("Err", errors.New("problem in uploaded file"), err)
			showError(c, err)
			return
		}

		filename := header.Filename

		out, err := os.Create("/tmp/" + filename)
		if err != nil {
			plog.Error("Err", errors.New("Cannot create file"), err)
			showError(c, err)
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			plog.Error("Err", errors.New("Cannot copy file"), err)
			showError(c, err)
			return
		}

		plog.Info("Sent file for upload: ", filename)

		err = UploadCsv(filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	}
}

func showError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
