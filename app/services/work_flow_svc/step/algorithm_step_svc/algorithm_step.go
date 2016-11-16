package algorithm_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/clients"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
	"github.com/crowdflux/angel/app/services/work_flow_io_svc"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
)

type algorithmStep struct {
	step.Step
	stepConfigSvc work_flow_io_svc.IStepConfigSvc
}

func (t *algorithmStep) processFlu(flu feed_line.FLU) {
	t.AddToBuffer(flu)

	plog.Info("algorithm Step flu reached", flu.ID)
	tStep, err := t.stepConfigSvc.GetAlgorithmStepConfig(flu.StepId)
	if err != nil {
		plog.Error("Algorithm step", err, "fluId: "+flu.ID.String(), "stepid: "+flu.StepId.String(), flu.FeedLineUnit)
		flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.Algorithm, "Algorithm Config Error", flu.Redelivered())
		t.finishFlu(flu)
		flu.ConfirmReceive()
		return
	}

	text := flu.Data[tStep.TextFieldKey]

	textSting := text.(string)

	algoResult, err, success := clients.GetAbacusClient().Predict(textSting)
	if err != nil {
		plog.Error("Algorithm step", err, "fluId: "+flu.ID.String(), flu.FeedLineUnit)
		flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.Algorithm, "Algorithm Error", flu.Redelivered())
	} else if success {
		flu.Build.Merge(models.JsonF{tStep.AnswerKey: algoResult})
	}

	t.finishFlu(flu)
	flu.ConfirmReceive()
}

func (t *algorithmStep) finishFlu(flu feed_line.FLU) bool {

	err := t.RemoveFromBuffer(flu)
	if err != nil {
		plog.Trace("algorithm step", "flu not present in buffer")
	}
	t.OutQ.Push(flu)
	plog.Info("algorithm step out", flu.ID)
	flu_logger_svc.LogStepExit(flu.FeedLineUnit, step_type.Algorithm, flu.Redelivered())

	return true
}

func newStdPredictor() *algorithmStep {
	ts := &algorithmStep{
		Step:          step.New(step_type.Algorithm),
		stepConfigSvc: work_flow_io_svc.NewStepConfigService(),
	}

	ts.SetFluProcessor(ts.processFlu)
	return ts
}
