package handlers

import (
	"fmt"
	"net/http"

	"errors"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/services/image_svc/downloader"
)

func BulkDownloadImages(c *gin.Context) {

	var images []models.ImageContainer

	if r := c.BindJSON(images); r != nil {
		fmt.Println("Error ouccured: ", r)
		return
	}
	downloader.DownloadFromList(images, "test")

	c.JSON(http.StatusOK, images)

}
func BulkDownloadedImagesFromCSV(c *gin.Context) {

	file, header, err := c.Request.FormFile("upload")
	if err != nil {
		showError(c, errors.New("problem in uploaded file"))
		return
	}

	filename := header.Filename

	folderName := c.PostForm("folderName")
	if folderName == "" {
		showError(c, errors.New("provide folderName"))
		return
	}

	if _, err := os.Stat("./temp"); os.IsNotExist(err) {
		if err = os.Mkdir("./temp", os.ModePerm); err != nil {
			showError(c, err)
			return
		}
	}

	out, err := os.Create("./temp/" + filename)
	if err != nil {
		showError(c, err)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		showError(c, err)
		return
	}

	images, err := downloader.ReadFromTempDir(filename)
	if err != nil {
		showError(c, err)
		return
	}

	err = downloader.DownloadFromList(images, folderName)
	if err != nil {
		showError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func showError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
