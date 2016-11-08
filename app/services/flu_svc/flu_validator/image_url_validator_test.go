package flu_validator

import (
	"github.com/asaskevich/govalidator"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

var fluId = uuid.NewV4()
var image_url_valid = []string{"https://s3-ap-southeast-1.amazonaws.com/playmentproduction/public/B00X0X3AKG_2.jpg", "https://s3-ap-southeast-1.amazonaws.com/playmentproduction/public/B00PU0DELW_2.jpg"}
var image_url_invalid = []string{"https://s-ap-southeast-1.amazonaws.com/playmentproduction/public/B00X0X3AKG_2.jpg", "https://s3-ap-southeast-1.amazonaws.com/playmentproduction/public/B00PU0DELW_2.jpg"}

var valid_flu = models.FeedLineUnit{
	ID:          fluId,
	ReferenceId: "PayFlip123",
	Tag:         "Ola",
	Data: models.JsonF{
		"brand":     "Sony",
		"image_url": image_url_valid},
	Build: models.JsonF{},
}

var invalid_url_flu = models.FeedLineUnit{
	ID:          fluId,
	ReferenceId: "PayFlip123",
	Tag:         "Ola",
	Data: models.JsonF{
		"brand":     "Sony",
		"image_url": image_url_invalid},
	Build: models.JsonF{},
}

var input_config = []models.FLUValidator{
	{uuid.NewV4(), "image_url", "image_array", false, "Ola", pq.NullTime{}, pq.NullTime{}, uuid.NewV4()},
}

var invalid_input_config = []models.FLUValidator{
	{uuid.NewV4(), "image_ur", "image_array", false, "Ola", pq.NullTime{}, pq.NullTime{}, uuid.NewV4()},
}

func Test_for_valid_urls(t *testing.T) {
	initialUrl := valid_flu.Data["image_url"]
	err := imageUrlEncryptor(&valid_flu, input_config)
	returnedUrlList := valid_flu.Data["image_url"].([]string)
	assert.True(t, govalidator.IsURL(returnedUrlList[0]))
	assert.True(t, govalidator.IsURL(returnedUrlList[1]))
	assert.Nil(t, err)
	assert.EqualValues(t, len(returnedUrlList), 2)
	assert.NotEqual(t, valid_flu.Data["image_url"], initialUrl)
}

func Test_for_invalid_urls(t *testing.T) {
	err := imageUrlEncryptor(&invalid_url_flu, input_config)
	assert.Error(t, err)
}

func Test_for_invalid_config(t *testing.T) {
	err := imageUrlEncryptor(&valid_flu, invalid_input_config)
	assert.Error(t, err)
}
