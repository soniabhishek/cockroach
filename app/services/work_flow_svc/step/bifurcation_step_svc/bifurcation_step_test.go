package bifurcation_step_svc

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
	config.Multiplication = 4
	return
}
func (s *stepConfigSvcMock) GetUnificationStepConfig(stepId uuid.UUID) (config models.UnificationConfig, err error) {
	config.Multiplication = 4
	return
}

func TestBifurcation_ProcessFlu(t *testing.T) {

	bfs := bifurcationStep{
		Step:          step.New(step_type.Test),
		stepConfigSvc: &stepConfigSvcMock{},
		fluRepo:       feed_line_repo.Mock(),
	}

	bfs.Step.SetFluProcessor(bfs.processFlu)

	bfs.Start()

	fluId := uuid.NewV4()
	inputFlu := feed_line.FLU{
		FeedLineUnit: models.FeedLineUnit{
			ID:       fluId,
			IsMaster: true,
			MasterId: fluId,
		},
	}

	bfs.InQ.Push(inputFlu)
	i := 0

	for flu := range bfs.OutQ.Receiver() {
		flu.ConfirmReceive()

		if i == 0 {
			assert.EqualValues(t, inputFlu.ID, flu.ID)
			assert.True(t, flu.IsMaster)
		} else {
			assert.EqualValues(t, inputFlu.MasterId, flu.MasterId)
			assert.False(t, flu.IsMaster)
		}

		assert.EqualValues(t, flu.Build[index], i)
		i++

		if i >= 4 {
			break
		}
	}

	assert.EqualValues(t, 4, i)

}
