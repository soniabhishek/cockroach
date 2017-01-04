package flu_monitor

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
	"time"
)

func TestFluMonitor_AddToOutputQueue(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	requestCountP1 := 0
	requestCountP2 := 0

	httpmock.RegisterResponder("POST", "https://url.somedomain.com/hit", func(req *http.Request) (*http.Response, error) {
		requestCountP1 = requestCountP1 + 1
		return httpmock.NewStringResponse(500, `{}`), nil
	},
	)

	httpmock.RegisterResponder("POST", "https://url.someotherdomain.com/hit", func(req *http.Request) (*http.Response, error) {
		requestCountP2 = requestCountP2 + 1
		return httpmock.NewStringResponse(200, `{"invalid_flus": []}`), nil
	},
	)

	fm := New()
	sendFlus(*fm)

	time.Sleep(time.Duration(25) * time.Second)
	//project with flu_count limit 1
	assert.Equal(t, 60, requestCountP1)
	//project with no flu_count limit per request
	assert.Equal(t, 1, requestCountP2)
}

func sendFlus(fm FluMonitor) {
	pID1 := uuid.FromStringOrNil("38a27427-9474-4e76-bd77-bc7dee2463e5")
	pID2 := uuid.FromStringOrNil("afa0c4f4-4998-4412-b3f2-039ae32b13ab")

	go func() {
		for i := 0; i < 20; i++ {
			flu := models.FeedLineUnit{
				ID:          uuid.NewV4(),
				ReferenceId: "38a27427-9474-4e76-bd77-bc7dee2463e5",
				ProjectId:   pID1,
				Data: models.JsonF{
					"product_id": "40843808",
				},
				Build: models.JsonF{
					"product_id": "40843808",
					"result":     "success",
				},
				Tag: "PAYTM_TSHIRT",
			}
			fm.AddToOutputQueue(flu)
		}
	}()
	go func() {
		for i := 0; i < 40; i++ {
			flu := models.FeedLineUnit{
				ID:          uuid.NewV4(),
				ReferenceId: "38a27427-9474-4e76-bd77-bc7dee2463e5",
				ProjectId:   pID2,
				Data: models.JsonF{
					"product_id": "40843808",
				},
				Build: models.JsonF{
					"product_id": "40843808",
					"result":     "success",
				},
				Tag: "PAYTM_TSHIRT",
			}
			fm.AddToOutputQueue(flu)
		}
	}()

}
