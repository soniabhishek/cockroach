package unification_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
	"github.com/crowdflux/angel/app/services/work_flow_io_svc"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
)

type unificationStep struct {
	step.Step
	stepConfigSvc work_flow_io_svc.IStepConfigurationSvc

	fluCounter fluCounter
}

func (u *unificationStep) processFlu(flu feed_line.FLU) {

	unificationConfig, err := u.stepConfigSvc.GetUnificationStepConfig(flu.StepId)
	if err != nil {
		plog.Error("Bifurcation Step", err, "error getting step", "fluId: "+flu.ID.String())
		flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.Unification, "Error getting bifurcation Config", flu.Redelivered())
		return
	}

	u.fluCounter.UpdateCount(flu)

	// If count is less then awaited count
	// then add it to the counter & don't finish it

	if u.fluCounter.GetCount(flu.ID) >= unificationConfig.Multiplication {

		waitingFlus := u.fluCounter.Get(flu.ID)

		for _, wFlu := range waitingFlus {
			flu.Build.Merge(wFlu.Build)
		}

		u.fluCounter.Clear(flu.ID)
		u.finishFlu(flu)
	}

	flu.ConfirmReceive()
}

func (u *unificationStep) finishFlu(flu feed_line.FLU) bool {
	u.OutQ.Push(flu)
	flu_logger_svc.LogStepExit(flu.FeedLineUnit, step_type.Unification, flu.Redelivered())
	return true
}
