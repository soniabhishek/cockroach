package logic_gate_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
	"github.com/crowdflux/angel/app/services/work_flow_svc/counter"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
)

type logic_gate_Step struct {
	step.Step
}

func (m *logic_gate_Step) processFlu(flu feed_line.FLU) {
	flu.ConfirmReceive()
	plog.Info("Manual Step flu reached", flu.ID)
	counter.Print(flu, "manual")
	m.finishFlu(flu)
}

func (m *logic_gate_Step) finishFlu(flu feed_line.FLU) bool {
	m.OutQ.Push(flu)
	flu_logger_svc.LogStepExit(flu.FeedLineUnit, step_type.Manual, flu.Redelivered())
	return true
}

func newLogicGateStep() *logic_gate_Step {
	ms := &logic_gate_Step{
		Step: step.New(step_type.LogicGate),
	}
	ms.Step.SetFluProcessor(ms.processFlu)
	return ms
}
