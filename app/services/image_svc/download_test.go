package image_svc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadImageSuccess(t *testing.T) {
	t.SkipNow()

	//Please change the image link if it breaks or expires
	googleLogoUrl := "https://www.google.co.in/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"

	c := make(chan downloadedImage, 1)

	download2(googleLogoUrl, c)
	image := <-c

	assert.NoError(t, image.err, "Image Download Failed")
}

func TestDownloadImageFail(t *testing.T) {

	t.SkipNow()

	//Some invalid link
	brokenUrl := "ht1233213tps://www.google"

	c := make(chan downloadedImage, 1)

	download2(brokenUrl, c)
	image := <-c

	assert.Error(t, image.err, "Image Download Succeeded")
}

func TestUploadImage(t *testing.T) {

}

func TestImageResize(t *testing.T) {

}
