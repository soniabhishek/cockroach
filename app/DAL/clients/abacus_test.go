package clients

import (
	"encoding/json"
	"github.com/crowdflux/angel/app/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAbacusPrediction(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var request algorithmRequest
	var err error
	httpmock.RegisterResponder("POST", config.ABACUS_API.Get()+"/api/review", func(req *http.Request) (*http.Response, error) {
		body, _ := ioutil.ReadAll(req.Body)
		err = json.Unmarshal(body, &request)
		if request.Input == "Good item" {
			return httpmock.NewStringResponse(200, `{"prediction":"Approve", "success":true}`), nil
		} else {
			return httpmock.NewStringResponse(200, `{ "success":false}`), nil
		}
	},
	)

	cc := GetAbacusClient()

	actualGoodResult, err, success := cc.Predict("Good item")
	expectedGoodResult := "Approve"

	assert.Equal(t, expectedGoodResult, actualGoodResult)
	assert.True(t, success)
	assert.NoError(t, err)

	actualBadResult, err, success := cc.Predict("Jango Fett")
	assert.NotEqual(t, actualBadResult, "Approve")
	assert.False(t, success)
	assert.NoError(t, err)
}
