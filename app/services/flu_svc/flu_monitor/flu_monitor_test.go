package flu_monitor

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	//	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"testing"
	"time"
)

func TestFluMonitor_AddToOutputQueue(t *testing.T) {
	flu := models.FeedLineUnit{
		ID:          uuid.FromStringOrNil("18846509-8ec6-49f6-8312-7d5d176fc0ac"),
		ReferenceId: "38a27427-9474-4e76-bd77-bc7dee2463e5",
		ProjectId:   uuid.FromStringOrNil("38a27427-9474-4e76-bd77-bc7dee2463e5"),
		Data: models.JsonF{
			"product_id": "40843808",
			"result":     "success",
		},
		Build: models.JsonF{
			"product_id": "40843808",
			"result":     "success",
		},
		Tag: "PAYTM_TSHIRT",
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://url.somedomain.com/hit",
		httpmock.NewStringResponder(200, `[{"id": 1, "invalid_flus": []`))

	fm := New()
	fm.AddToOutputQueue(flu)

	time.Sleep(time.Duration(20) * time.Second)
	//assert.NoError(t, err, "Error occured")
	//	assert.True(t, isValid, "Expected valid flu but found inValid")
}
