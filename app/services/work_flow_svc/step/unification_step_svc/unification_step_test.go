package unification_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/work_flow_io_svc"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
	"github.com/stretchr/testify/assert"
	"testing"
)

type stepConfigSvcMock struct{}

var _ work_flow_io_svc.IStepConfigSvc = &stepConfigSvcMock{}

func (s *stepConfigSvcMock) GetTransformationStepConfig(stepId uuid.UUID) (config models.TransformationConfig, err error) {
	return
}
func (s *stepConfigSvcMock) GetBifurcationStepConfig(stepId uuid.UUID) (config models.BifurcationConfig, err error) {
	config.Multiplication = 3
	return
}
func (s *stepConfigSvcMock) GetUnificationStepConfig(stepId uuid.UUID) (config models.UnificationConfig, err error) {
	config.Multiplication = 3
	return
}

func TestUnification_ProcessFlu(t *testing.T) {

	fluRepoMock := feed_line_repo.Mock()

	unifStp := unificationStep{
		Step:          step.New(step_type.Test),
		stepConfigSvc: &stepConfigSvcMock{},
		fluRepo:       fluRepoMock,
		fluCounter:    newFluCounter(),
	}

	unifStp.Step.SetFluProcessor(unifStp.processFlu)

	unifStp.Start()

	id := uuid.NewV4()
	inputFlu := feed_line.FLU{
		FeedLineUnit: models.FeedLineUnit{
			ID:       id,
			Build:    models.JsonF{"prop0": "a"},
			StepId:   uuid.NewV4(),
			IsActive: true,
			IsMaster: true,
			MasterId: id,
		},
	}

	inputFlu2 := inputFlu
	inputFlu2.ID = uuid.NewV4()
	inputFlu2.IsMaster = false
	inputFlu2.MasterId = inputFlu.MasterId
	inputFlu2.Build = models.JsonF{"prop1": 11}

	inputFlu3 := inputFlu
	inputFlu3.ID = uuid.NewV4()
	inputFlu3.IsMaster = false
	inputFlu3.MasterId = inputFlu.MasterId
	inputFlu3.Build = models.JsonF{"prop2": true}

	fluRepoMock.Add(inputFlu.FeedLineUnit)
	fluRepoMock.Add(inputFlu2.FeedLineUnit)
	fluRepoMock.Add(inputFlu3.FeedLineUnit)

	unifStp.InQ.Push(inputFlu)
	unifStp.InQ.Push(inputFlu2)
	unifStp.InQ.Push(inputFlu3)

	r := unifStp.OutQ.Receiver()

	flu := <-r

	flu.ConfirmReceive()

	assert.EqualValues(t, inputFlu.ID, flu.ID)
	assert.True(t, flu.IsMaster)
	assert.True(t, flu.IsActive)
	assert.Equal(t, "a", flu.Build["prop0"])
	assert.EqualValues(t, 11, flu.Build["prop1"])
	assert.Equal(t, true, flu.Build["prop2"])

	select {
	case flu = <-r:
		flu.ConfirmReceive()
		assert.Fail(t, "more value recieved")
	default:
	}
}
