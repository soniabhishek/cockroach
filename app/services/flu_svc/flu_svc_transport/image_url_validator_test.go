package flu_svc_transport_test

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_svc_transport"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

var fluId = uuid.NewV4()

var flu = models.FeedLineUnit{
	ID:          fluId,
	ReferenceId: "PayFlip123",
	Tag:         "Ola",
	Data: models.JsonF{
		"brand":     "Sony",
		"image_url": `["https://s3-ap-southeast-1.amazonaws.com/playmentproduction/public/B00X0X3AKG_2.jpg"]`,
	},
	Build: models.JsonF{},
}
var input_config = []models.FLUValidator{
	{uuid.NewV4(), "image", "image_url", false, "Ola", pq.NullTime{}, pq.NullTime{}, uuid.NewV4()},
}

func Test(t *testing.T) {
	flu_svc_transport.ImageUrlValidator(&flu, input_config)
	//assert.EqualValues(t, flu.ID, fluNew.ID)
	assert.EqualValues(t, flu.Build["new_prop"], 123)
}
