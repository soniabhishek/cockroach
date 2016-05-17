package image_svc1

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/support/app/models"
	"gitlab.com/playment-main/support/app/services/image_svc/downloader"
	"testing"
)

func TestImageService_BulkDownloadImages(t *testing.T) {
	imageSvc := &imageService{}

	images, err := downloader.ReadFromTempDir("images_myntra_top_pose.csv")

	if err != nil {
		assert.FailNow(t, "reading csv failed"+err.Error())
	}
	var imageDictionaries []models.ImageDictionary = make([]models.ImageDictionary, len(images))

	for i, v := range images {
		imageDictionaries[i] = models.ImageDictionary{
			Label:       v.Id,
			OriginalUrl: v.Url,
		}
	}

	_, err = imageSvc.BulkDownloadImages(imageDictionaries)
	assert.NoError(t, err)
}
