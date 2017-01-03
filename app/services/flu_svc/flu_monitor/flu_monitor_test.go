package flu_monitor

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	//	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"testing"
	"time"
)

var flu1 models.FeedLineUnit
var flu2 models.FeedLineUnit
var flu3 models.FeedLineUnit
var flu4 models.FeedLineUnit

func Setup() {
	fluID1 := uuid.NewV4()
	fluID2 := uuid.NewV4()
	fluID3 := uuid.NewV4()
	fluID4 := uuid.NewV4()
	pID1 := uuid.FromStringOrNil("38a27427-9474-4e76-bd77-bc7dee2463e5")
	pID2 := uuid.FromStringOrNil("afa0c4f4-4998-4412-b3f2-039ae32b13ab")
	flu1 = models.FeedLineUnit{
		ID:          fluID1,
		ReferenceId: "38a27427-9474-4e76-bd77-bc7dee2463e5",
		ProjectId:   pID1,
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

	flu2 = models.FeedLineUnit{
		ID:          fluID2,
		ReferenceId: "38a27427-9474-4e76-bd77-bc7dee2463e5",
		ProjectId:   pID1,
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
	flu3 = models.FeedLineUnit{
		ID:          fluID3,
		ReferenceId: "38a27427-9474-4e76-bd77-bc7dee2463e5",
		ProjectId:   pID2,
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
	flu4 = models.FeedLineUnit{
		ID:          fluID4,
		ReferenceId: "38a27427-9474-4e76-bd77-bc7dee2463e5",
		ProjectId:   pID2,
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
}

func TestFluMonitor_AddToOutputQueue(t *testing.T) {

	Setup()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://url.somedomain.com/hit",
		httpmock.NewStringResponder(200, `{"invalid_flus": []}`))

	httpmock.RegisterResponder("POST", "https://url.someotherdomain.com/hit",
		httpmock.NewStringResponder(200, `{"invalid_flus": []}`))

	fm := New()
	fm.AddToOutputQueue(flu1)
	fm.AddToOutputQueue(flu2)
	fm.AddToOutputQueue(flu3)
	fm.AddToOutputQueue(flu4)

	time.Sleep(time.Duration(20) * time.Second)
	//assert.NoError(t, err, "Error occured")
	//	assert.True(t, isValid, "Expected valid flu but found inValid")
}
