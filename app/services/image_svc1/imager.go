package image_svc1

import (
	"errors"
	"io"
	"net/http"

	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/services/data_access_svc/clients"
)

//Add Downloading and Uploading functionality to models.ImageDictionary
type imager struct {
	models.ImageDictionary
	imgData io.ReadCloser
}

//Downloads the image into an internal readCloser
func (i *imager) Download() error {

	if i.OriginalUrl == "" {
		errors.New("Original Url missing")
	}

	resp, err := http.Get(i.OriginalUrl)
	if err != nil {
		return err
	}
	i.imgData = resp.Body
	return nil
}

//Uploads the image to aws S3
func (i *imager) Upload() error {
	err := clients.GetS3Client().Upload(i.imgData, "playmentdevelopment", "default/test1/"+i.Label+".png")
	i.imgData.Close()
	return err
}

//Re-sizes the downloaded image....Not implemented yet
func (i *imager) Resize() error {
	return errors.New("Not implemented")
}
