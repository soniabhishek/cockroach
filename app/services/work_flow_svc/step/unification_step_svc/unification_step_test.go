package unification_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/work_flow_io_svc"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
	"github.com/stretchr/testify/assert"
	"testing"
)

type stepConfigSvcMock struct{}

var _ work_flow_io_svc.IStepConfigurationSvc = &stepConfigSvcMock{}

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

	unifStp := unificationStep{
		Step:          step.New(step_type.Test),
		stepConfigSvc: &stepConfigSvcMock{},
		fluCounter:    newFluCounter(),
	}

	unifStp.Step.SetFluProcessor(unifStp.processFlu)

	unifStp.Start()

	inputFlu := feed_line.FLU{
		FeedLineUnit: models.FeedLineUnit{
			ID:     uuid.NewV4(),
			Build:  models.JsonF{"prop0": "a"},
			CopyId: 0,
		},
	}

	inputFlu2 := inputFlu
	inputFlu2.Build = models.JsonF{"prop1": 11}
	inputFlu2.CopyId = 1

	inputFlu3 := inputFlu
	inputFlu3.Build = models.JsonF{"prop2": true}
	inputFlu3.CopyId = 2

	unifStp.InQ.Push(inputFlu)
	unifStp.InQ.Push(inputFlu2)
	unifStp.InQ.Push(inputFlu3)

	r := unifStp.OutQ.Receiver()

	flu := <-r

	flu.ConfirmReceive()

	assert.EqualValues(t, inputFlu.ID, flu.ID)
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
