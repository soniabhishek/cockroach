package start_step_svc

import (
	"testing"

	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/work_flow_io_svc"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
	"github.com/stretchr/testify/assert"
	"time"
)

//--------------------------------------------------------------------------------//

type stepConfigSvcMock struct{}

var _ work_flow_io_svc.IStepConfigSvc = &stepConfigSvcMock{}

type fakeFluPusher struct {
}

func (fakeFluPusher) PushFLU(models.FeedLineUnit) (bool, error) {
	return true, nil
}

var fluId = uuid.NewV4()

var flu = feed_line.FLU{
	FeedLineUnit: models.FeedLineUnit{
		ID:          fluId,
		ReferenceId: "PayFlip123",
		Tag:         "Brand",
		Data: models.JsonF{
			"brand":     "Otter",
			"image_url": "[\"https://s3-ap-southeast-1.amazonaws.com/playmentproduction/public/B00X0X3AKG_2.jpg\",\"https://s3-ap-southeast-1.amazonaws.com/playmentproduction/public/B00PU0DELW_2.jpg\"",
		},
		Build: models.JsonF{},
	},
}

func (s *stepConfigSvcMock) GetTransformationStepConfig(stepId uuid.UUID) (config models.TransformationConfig, err error) {
	return
}
func (s *stepConfigSvcMock) GetBifurcationStepConfig(stepId uuid.UUID) (config models.BifurcationConfig, err error) {
	//	config.Multiplication = 2
	return
}
func (s *stepConfigSvcMock) GetUnificationStepConfig(stepId uuid.UUID) (config models.UnificationConfig, err error) {
	//	config.Multiplication = 2
	return
}

func (s *stepConfigSvcMock) GetCrowdsourcingStepConfig(stepId uuid.UUID) (config models.CrowdsourcingConfig, err error) {
	return
}

func (s *stepConfigSvcMock) GetStartStepConfig(stepId uuid.UUID) (config models.StartStepConfig, err error) {
	config.ImageFieldKey = "image_url"
	return
}

func Test(t *testing.T) {

	fluRepo := feed_line_repo.Mock()

	fluRepo.Save(flu.FeedLineUnit)

	cs := starterStep{
		Step: step.New(step_type.Test),
	}

	cs.SetFluProcessor(cs.processFlu)

	cs.Start()

	defer cs.finishFlu(flu)

	cs.InQ.Push(flu)

	// Giving it time to finish adding to buffer
	// as its happening in another goroutine
	time.Sleep(time.Duration(100) * time.Millisecond)

	plog.Debug("fluData", flu.Data)

	flu.Build["new_prop"] = 123

	ok := cs.finishFlu(flu)
	assert.True(t, ok)

	var fluNew feed_line.FLU
	select {
	case fluNew = <-cs.OutQ.Receiver():
		fluNew.ConfirmReceive()
		assert.EqualValues(t, flu.ID, fluNew.ID)
		assert.EqualValues(t, flu.Build["new_prop"], 123)
	case <-time.After(time.Duration(2) * time.Second):
		assert.FailNow(t, "nothing came out of crowdsourcing queue")
	}

}
