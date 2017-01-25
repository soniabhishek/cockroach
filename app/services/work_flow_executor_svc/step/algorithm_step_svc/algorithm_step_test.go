package algorithm_step_svc

import (
	"testing"

	"encoding/json"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/work_flow_executor_svc/step"
	"github.com/crowdflux/angel/app/services/work_flow_svc"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"io/ioutil"
	"net/http"
	"time"
)

type stepConfigSvcMock struct{}

var _ work_flow_svc.IStepConfigSvc = &stepConfigSvcMock{}

type fakeFluPusher struct {
}

func (fakeFluPusher) PushFLU(models.FeedLineUnit) (bool, error) {
	return true, nil
}

var fluId = uuid.NewV4()
var badFluId = uuid.NewV4()

var flu = feed_line.FLU{
	FeedLineUnit: models.FeedLineUnit{
		ID:          fluId,
		ReferenceId: "PayFlip123",
		Tag:         "Brand",
		Data: models.JsonF{
			"review_body": "Good product",
		},
		Build: models.JsonF{
			"review_body": "Good product",
		},
	},
}

var badFlu = feed_line.FLU{
	FeedLineUnit: models.FeedLineUnit{
		ID:          badFluId,
		ReferenceId: "PayFlip123",
		Tag:         "Brand",
		Data: models.JsonF{
			"review_body": "Darth Vader",
		},
		Build: models.JsonF{
			"review_body": "Darth Vader",
			"Luke":        "I am your father",
		},
	},
}

func (s *stepConfigSvcMock) GetAlgorithmStepConfig(stepId uuid.UUID) (config models.AlgorithmConfig, err error) {
	config.AnswerKey = "algo_result"
	config.TextFieldKey = "review_body"
	return
}

func (s *stepConfigSvcMock) GetTransformationStepConfig(stepId uuid.UUID) (config models.TransformationConfig, err error) {
	return
}

func (s *stepConfigSvcMock) GetValidationStepConfig(stepId uuid.UUID) (config models.ValidationConfig, err error) {
	return
}

func (s *stepConfigSvcMock) GetBifurcationStepConfig(stepId uuid.UUID) (config models.BifurcationConfig, err error) {
	return
}
func (s *stepConfigSvcMock) GetUnificationStepConfig(stepId uuid.UUID) (config models.UnificationConfig, err error) {
	return
}

func (s *stepConfigSvcMock) GetCrowdsourcingStepConfig(stepId uuid.UUID) (config models.CrowdsourcingConfig, err error) {
	return
}

type algorithmRequest struct {
	Input string `json:"review"`
}

func TestSuccessfulPrediction(t *testing.T) {

	fluRepo := feed_line_repo.Mock()

	fluRepo.Save(flu.FeedLineUnit)

	as := algorithmStep{
		Step:          step.New(step_type.Test),
		stepConfigSvc: &stepConfigSvcMock{},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var request algorithmRequest
	var err error
	httpmock.RegisterResponder("POST", config.ABACUS_API.Get()+"/api/review", func(req *http.Request) (*http.Response, error) {
		body, _ := ioutil.ReadAll(req.Body)
		err = json.Unmarshal(body, &request)
		return httpmock.NewStringResponse(200, `{"prediction":"Approve", "success":true}`), nil
	},
	)

	as.SetFluProcessor(as.processFlu)

	as.Start()

	as.InQ.Push(flu)

	// Giving it time to finish adding to buffer
	// as its happening in another goroutine
	time.Sleep(time.Duration(100) * time.Millisecond)

	var fluNew feed_line.FLU
	select {
	case fluNew = <-as.OutQ.Receiver():
		fluNew.ConfirmReceive()
		assert.EqualValues(t, flu.ID, fluNew.ID)
		assert.EqualValues(t, "Good product", request.Input)
		assert.EqualValues(t, "Approve", fluNew.Build["algo_result"])
	case <-time.After(time.Duration(4) * time.Second):
		assert.FailNow(t, "nothing came out of crowdsourcing queue")
	}

}

func TestUnSuccessfulPrediction(t *testing.T) {

	fluRepo := feed_line_repo.Mock()

	fluRepo.Save(badFlu.FeedLineUnit)

	cs := algorithmStep{
		Step:          step.New(step_type.Test),
		stepConfigSvc: &stepConfigSvcMock{},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var request algorithmRequest
	var err error
	httpmock.RegisterResponder("POST", config.ABACUS_API.Get()+"/api/review", func(req *http.Request) (*http.Response, error) {
		body, _ := ioutil.ReadAll(req.Body)
		err = json.Unmarshal(body, &request)
		return httpmock.NewStringResponse(200, `{ "success":false}`), nil
	},
	)
	cs.SetFluProcessor(cs.processFlu)

	cs.Start()

	cs.InQ.Push(badFlu)

	// Giving it time to finish adding to buffer
	// as its happening in another goroutine
	time.Sleep(time.Duration(100) * time.Millisecond)

	var fluNew feed_line.FLU
	select {
	case fluNew = <-cs.OutQ.Receiver():
		fluNew.ConfirmReceive()
		assert.EqualValues(t, badFlu.ID, fluNew.ID)
		assert.EqualValues(t, badFlu.Build.Merge(models.JsonF{"algo_result_success": false}), fluNew.Build)
	case <-time.After(time.Duration(2) * time.Second):
		assert.FailNow(t, "nothing came out of crowdsourcing queue")
	}

}
