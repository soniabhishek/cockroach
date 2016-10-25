package bifurcation_step_svc

import (
	"fmt"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
	"github.com/crowdflux/angel/app/services/work_flow_io_svc"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
)

type bifurcationStep struct {
	step.Step
	stepConfigSvc work_flow_io_svc.IStepConfigurationSvc
}

func (b *bifurcationStep) processFlu(flu feed_line.FLU) {

	bifurcationConfig, err := b.stepConfigSvc.GetBifurcationStepConfig(flu.StepId)
	if err != nil {
		plog.Error("Bifurcation Step", err, "error getting step", "fluId: "+flu.ID.String())
		flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.Bifurcation, "Error getting bifurcation Config", flu.Redelivered())
		return
	}

	if flu.Build == nil {
		flu.Build = models.JsonF{}
	}

	if bifurcationConfig.Multiplication > 0 {

		for i := 0; i < bifurcationConfig.Multiplication; i++ {

			newFlu := flu
			newFlu.CopyId = i

			b.finishFlu(newFlu)
		}

	} else {
		plog.Error("Bifurcation Step", fmt.Errorf("", "Multiplication count not greater than 0"), "flu_id "+flu.ID.String())
		b.finishFlu(flu)
	}

	flu.ConfirmReceive()
}

func (c *bifurcationStep) finishFlu(flu feed_line.FLU) bool {
	c.OutQ.Push(flu)
	flu_logger_svc.LogStepExit(flu.FeedLineUnit, step_type.Bifurcation, flu.Redelivered())
	return true
}
