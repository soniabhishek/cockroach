package image_svc1

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"

	"testing"
)

func TestImager_Download(t *testing.T) {

	t.SkipNow()

	imageId := uuid.NewV4()
	imageUrl := "https://pixabay.com/static/uploads/photo/2015/10/01/21/39/background-image-967820_960_720.jpg"
	brokenImageUrl := "https://pixabay.20_9.jpg/"

	imagerObj := &imager{
		ImageDictionaryNew: models.ImageDictionaryNew{
			ID:          imageId,
			OriginalUrl: imageUrl,
		},
	}

	err := imagerObj.Download()
	assert.NoError(t, err, "Should not throw error")

	imagerObj = &imager{
		ImageDictionaryNew: models.ImageDictionaryNew{
			ID:          imageId,
			OriginalUrl: brokenImageUrl,
		},
	}

	err = imagerObj.Download()
	assert.Error(t, err)
}

func TestImager_Upload(t *testing.T) {

	t.SkipNow()

	imageId := uuid.NewV4()
	imageUrl := "https://pixabay.com/static/uploads/photo/2015/10/01/21/39/background-image-967820_960_720.jpg"

	imagerObj := &imager{
		ImageDictionaryNew: models.ImageDictionaryNew{
			ID:          imageId,
			OriginalUrl: imageUrl,
		},
	}

	err := imagerObj.Download()
	assert.NoError(t, err)

	err = imagerObj.Upload()
	assert.NoError(t, err)
}
