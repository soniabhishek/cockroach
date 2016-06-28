package manual_step

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
	"gitlab.com/playment-main/angel/utilities/constants"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

var StdManualStep = &manualStep{
	Step: step.New(),
}

// Just a short form for above
var Std = StdManualStep

func AddHttpTransport(r *gin.RouterGroup) {

	r.POST(DOWNLOAD_ENDPOINT, fileDownloadHandler())
	r.POST(UPLOAD_ENDPOINT, fileUploadHandler())
}

func fileDownloadHandler() gin.HandlerFunc {

	return func(c *gin.Context) {

		manualStepId, err := uuid.FromString(c.PostForm(MANUAL_STEP_ID))
		plog.Info(c.PostForm(MANUAL_STEP_ID), manualStepId, err)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				SUCCESS: false,
				ERROR:   "please provide params",
			})
			return
		}

		fileUrl, err := DownloadCsv(manualStepId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				SUCCESS: false,
				ERROR:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			SUCCESS: true,
			//FILEPATH: "/downloadedfiles" + string(os.PathSeparator) + manualStepId.String() + ".csv",
			FILEPATH: fileUrl,
		})
		plog.Info(c.Request.RequestURI, fileUrl)
	}
}
func fileUploadHandler() gin.HandlerFunc {

	return func(c *gin.Context) {

		file, header, err := c.Request.FormFile(UPLOAD)
		if err != nil {
			plog.Error("Err", errors.New("problem in uploaded file"), err)
			showError(c, err)
			return
		}

		filename := header.Filename

		out, err := os.Create(TEMP_FOLDER + filename)
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
				SUCCESS: false,
				ERROR:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			SUCCESS: true,
		})
	}
}

func showError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		SUCCESS: false,
		ERROR:   err.Error(),
	})
}

//TODO send it to Utilities
func FlattenCSV(file string, url string, manualStepId uuid.UUID) (fileUrl string, err error) {

	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your image file
	f, err := os.Open(file)
	if err != nil {
		plog.Error("Error", err)
		return constants.Empty, err
	}
	defer f.Close()
	fw, err := w.CreateFormFile(PARAM_FILES, file)
	if err != nil {
		plog.Error("Error", err)
		return constants.Empty, err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}
	// Add the other fields
	//TODO check these, are they needed
	if fw, err = w.CreateFormField("key"); err != nil {
		plog.Error("Error", err)
		return
	}
	// TODO this one too.
	if _, err = fw.Write([]byte("KEY")); err != nil {
		plog.Error("Error", err)
		return
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		plog.Error("Error", err)
		return constants.Empty, err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set(CONTENT_TYPE, w.FormDataContentType())
	req.Header.Set(PARAM_PLAYMENT_ID, manualStepId.String())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		plog.Error("Error", err)
		return constants.Empty, err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		return
	}
	fmt.Println(res.StatusCode)

	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		plog.Error("Error", err)
		return constants.Empty, err
	}

	fileMap := models.JsonFake{}
	fileMap.Scan(string(response))
	plog.Info("FileMAP: ", fileMap)

	errTag := fileMap[ERROR]
	if errTag != nil {
		err = errors.New(strconv.FormatBool(errTag.(bool)))
		return
	}

	urlTag := fileMap[URL]
	if urlTag == nil {
		err = errors.New("No file url from Transformation.")
		return
	}
	fileUrl = urlTag.(string)
	plog.Info("FileURL: ", fileUrl)

	return
}
