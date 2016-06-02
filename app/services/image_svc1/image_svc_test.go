package image_svc1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/services/image_svc/downloader"
)

func TestImageService_BulkDownloadImages(t *testing.T) {
	t.SkipNow()

	imageSvc := &imageService{}

	images, err := downloader.ReadFromTempDir("images_myntra_top_pose.csv")

	if err != nil {
		assert.FailNow(t, "reading csv failed"+err.Error())
	}
	var imageDictionaries []models.ImageDictionaryNew = make([]models.ImageDictionaryNew, len(images))

	for i, v := range images {
		imageDictionaries[i] = models.ImageDictionaryNew{
			Label:       v.Id,
			OriginalUrl: v.Url,
		}
	}

	_, err = imageSvc.BulkDownloadImages(imageDictionaries)
	assert.NoError(t, err)
}
