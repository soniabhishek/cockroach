package image_svc1

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"

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
