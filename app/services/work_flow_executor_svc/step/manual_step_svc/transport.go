package manual_step_svc

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/utilities/constants"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

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
			SUCCESS:  true,
			FILEPATH: fileUrl,
		})
		plog.Info(c.Request.RequestURI, fileUrl)
	}
}
func fileUploadHandler() gin.HandlerFunc {

	return func(c *gin.Context) {

		plog.Trace("MANUAL", "Upload request reached")

		file, header, err := c.Request.FormFile(UPLOAD)
		if err != nil {
			plog.Error("Manual Step", err, plog.Message("problem in uploaded file"))
			showError(c, err)
			return
		}
		defer file.Close()

		filename := header.Filename
		plog.Trace("MANUAL", "File reached")
		out, err := os.Create(TEMP_FOLDER + filename)
		if err != nil {
			plog.Error("Manual Step", err, plog.Message("Cannot create file"))
			showError(c, err)
			return
		}
		defer out.Close()
		plog.Trace("COPY", "File copy is starting")
		_, err = io.Copy(out, file)
		if err != nil {
			plog.Error("Manual Step", err, plog.Message("Cannot copy file"))
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
		plog.Error("Manual Step", err, plog.Message("Flatten Csv. Error while opening file"))
		return constants.Empty, err
	}
	defer f.Close()
	fw, err := w.CreateFormFile(PARAM_FILES, file)
	if err != nil {
		plog.Error("Manual Step", err, plog.Message("Flatten Csv. Error while creating form file"))
		return constants.Empty, err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}
	// Add the other fields
	//TODO check these, are they needed
	if fw, err = w.CreateFormField("key"); err != nil {
		plog.Error("Manual Step", err, plog.Message("Flatten Csv. Error while creating form field"))
		return
	}
	// TODO this one too.
	if _, err = fw.Write([]byte("KEY")); err != nil {
		plog.Error("Manual Step", err, plog.Message("Flatten Csv. Error while writing to file"))
		return
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		plog.Error("Manual Step", err, plog.Message("Flatten Csv. Error in Post request"))
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
		plog.Error("Manual Step", err, plog.Message("Flatten Csv. Error in reading file"))
		return constants.Empty, err
	}

	fileMap := models.JsonF{}
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
