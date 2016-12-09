package empty_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
	"github.com/crowdflux/angel/app/services/work_flow_executor_svc/counter"
	"github.com/crowdflux/angel/app/services/work_flow_executor_svc/step"
)

type emptyStep struct {
	step.Step
}

func (m *emptyStep) processFlu(flu feed_line.FLU) {
	flu.ConfirmReceive()
	plog.Info("Empty step flu reached", flu.ID)
	counter.Print(flu, "empty")
	m.finishFlu(flu)
}

func (m *emptyStep) finishFlu(flu feed_line.FLU) bool {
	m.OutQ.Push(flu)
	flu_logger_svc.LogStepExit(flu.FeedLineUnit, step_type.EmptyStep, flu.Redelivered())
	return true
}

func newEmptyStep() *emptyStep {
	ms := &emptyStep{
		Step: step.New(step_type.EmptyStep),
	}
	ms.Step.SetFluProcessor(ms.processFlu)
	return ms
}
