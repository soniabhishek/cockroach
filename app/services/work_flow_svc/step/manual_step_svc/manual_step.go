package manual_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
	"github.com/crowdflux/angel/app/services/work_flow_svc/counter"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
)

type manualStep struct {
	step.Step
}

func (m *manualStep) processFlu(flu feed_line.FLU) {
	m.AddToBuffer(flu)
	flu.ConfirmReceive()
	plog.Info("Manual Step flu reached", flu.ID)
}

func (m *manualStep) finishFlu(flu feed_line.FLU) bool {

	err := m.RemoveFromBuffer(flu)
	if err != nil {
		plog.Trace("Manual Step", err, "flu not present", flu.ID)
		//return false
	}
	counter.Print(flu, "manual")
	m.OutQ.Push(flu)
	flu_logger_svc.LogStepExit(flu.FeedLineUnit, step_type.Manual)

	return true
}

func newManualStep() *manualStep {
	ms := &manualStep{
		Step: step.New(step_type.Manual),
	}
	ms.Step.SetFluProcessor(ms.processFlu)
	return ms
}
