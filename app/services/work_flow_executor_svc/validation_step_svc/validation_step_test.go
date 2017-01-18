package validation_step_svc

import (
	"encoding/json"
	"github.com/crowdflux/angel/app/DAL/feed_line"
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
	"testing"
	"time"
)

type stepConfigSvcMock struct{}

var _ work_flow_svc.IStepConfigSvc = &stepConfigSvcMock{}

type validationRequest struct {
	TemplateId string       `json:"template_id"`
	Input      models.JsonF `json:"input"`
}

func (*stepConfigSvcMock) GetAlgorithmStepConfig(stepId uuid.UUID) (config models.AlgorithmConfig, err error) {
	return
}

func (*stepConfigSvcMock) GetCrowdsourcingStepConfig(stepId uuid.UUID) (config models.CrowdsourcingConfig, err error) {
	return
}

func (s *stepConfigSvcMock) GetValidationStepConfig(stepId uuid.UUID) (config models.ValidationConfig, err error) {
	config.AnswerKey = "validation_result"
	config.TemplateId = "123"
	return
}

func (s *stepConfigSvcMock) GetTransformationStepConfig(stepId uuid.UUID) (config models.TransformationConfig, err error) {
	return
}
func (s *stepConfigSvcMock) GetBifurcationStepConfig(stepId uuid.UUID) (config models.BifurcationConfig, err error) {
	return
}
func (s *stepConfigSvcMock) GetUnificationStepConfig(stepId uuid.UUID) (config models.UnificationConfig, err error) {
	return
}

func TestValidation_ProcessFlu(t *testing.T) {

	vStep := validationStep{
		Step:          step.New(step_type.Test),
		stepConfigSvc: &stepConfigSvcMock{},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var request validationRequest
	var err error
	httpmock.RegisterResponder("POST", config.MEGATRON_API.Get()+"/validate", func(req *http.Request) (*http.Response, error) {
		body, _ := ioutil.ReadAll(req.Body)
		err = json.Unmarshal(body, &request)
		return httpmock.NewStringResponse(200, `{"success":true}`), nil
	},
	)

	vStep.Step.SetFluProcessor(vStep.processFlu)

	vStep.Start()

	id := uuid.NewV4()
	inputFlu := feed_line.FLU{
		FeedLineUnit: models.FeedLineUnit{
			ID:     id,
			Build:  models.JsonF{"prop0": "a"},
			StepId: uuid.NewV4(),
		},
	}
	vStep.InQ.Push(inputFlu)

	// Giving it time to finish adding to buffer
	// as its happening in another goroutine
	time.Sleep(time.Duration(100) * time.Millisecond)

	var outputFlu feed_line.FLU
	select {
	case outputFlu = <-vStep.OutQ.Receiver():
		outputFlu.ConfirmReceive()
		assert.EqualValues(t, inputFlu.ID, outputFlu.ID)
		assert.EqualValues(t, true, outputFlu.Build["validation_result"])
		assert.NoError(t, err)
		assert.Equal(t, request, validationRequest{TemplateId: "123", Input: models.JsonF{"prop0": "a"}})
	case <-time.After(time.Duration(2) * time.Second):
		assert.FailNow(t, "nothing came out of queue")
	}

}
