package algorithm_step_svc

import (
	"testing"

	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/work_flow_executor_svc/step"
	"github.com/crowdflux/angel/app/services/work_flow_io_svc"
	"github.com/stretchr/testify/assert"
	"time"
)

type stepConfigSvcMock struct{}

var _ work_flow_io_svc.IStepConfigSvc = &stepConfigSvcMock{}

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
		Build: models.JsonF{},
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
		Build: models.JsonF{"Luke": "I am your father"},
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
func (s *stepConfigSvcMock) GetBifurcationStepConfig(stepId uuid.UUID) (config models.BifurcationConfig, err error) {
	return
}
func (s *stepConfigSvcMock) GetUnificationStepConfig(stepId uuid.UUID) (config models.UnificationConfig, err error) {
	return
}

func (s *stepConfigSvcMock) GetCrowdsourcingStepConfig(stepId uuid.UUID) (config models.CrowdsourcingConfig, err error) {
	return
}
func TestSuccessfulPrediction(t *testing.T) {

	fluRepo := feed_line_repo.Mock()

	fluRepo.Save(flu.FeedLineUnit)

	cs := algorithmStep{
		Step:          step.New(step_type.Test),
		stepConfigSvc: &stepConfigSvcMock{},
	}

	cs.SetFluProcessor(cs.processFlu)

	cs.Start()

	cs.InQ.Push(flu)

	// Giving it time to finish adding to buffer
	// as its happening in another goroutine
	time.Sleep(time.Duration(100) * time.Millisecond)

	var fluNew feed_line.FLU
	select {
	case fluNew = <-cs.OutQ.Receiver():
		fluNew.ConfirmReceive()
		assert.EqualValues(t, flu.ID, fluNew.ID)
		assert.EqualValues(t, "Approve", fluNew.Build["algo_result"])
	case <-time.After(time.Duration(2) * time.Second):
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
		assert.EqualValues(t, badFlu.Build, fluNew.Build)
	case <-time.After(time.Duration(2) * time.Second):
		assert.FailNow(t, "nothing came out of crowdsourcing queue")
	}

}
