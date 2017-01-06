package transformation_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/clients"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
	"github.com/crowdflux/angel/app/services/work_flow_executor_svc/step"
	"github.com/crowdflux/angel/app/services/work_flow_svc"
)

type transformationStep struct {
	step.Step
	stepConfigSvc work_flow_svc.IStepConfigSvc
}

func (t *transformationStep) processFlu(flu feed_line.FLU) {
	t.AddToBuffer(flu)
	plog.Info("transformation Step flu reached", flu.ID)

	tStep, err := t.stepConfigSvc.GetTransformationStepConfig(flu.StepId)
	if err != nil {
		plog.Error("transformation step", err, plog.NewMessageWithParam("fluId: ", flu.ID.String()), plog.NewMessageWithParam("stepid: ", flu.StepId.String()), plog.NewMessageWithParam("flu", flu.FeedLineUnit))
		flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.Transformation, "TransformationConfigError", flu.Redelivered())
		return
	}

	transformedBuild, err := clients.GetMegatronClient().Transform(flu.Build, tStep.TemplateId)
	if err != nil {
		plog.Error("Transformation step", err, plog.NewMessageWithParam("fluId: ", flu.ID.String()), plog.NewMessageWithParam("flu", flu.FeedLineUnit))
		flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.Transformation, "TransformationError", flu.Redelivered())
		return
	}

	flu.Build.Merge(transformedBuild)

	t.finishFlu(flu)

}

func (t *transformationStep) finishFlu(flu feed_line.FLU) bool {

	err := t.RemoveFromBuffer(flu)
	if err != nil {
		plog.Trace("transformation step", "flu not present in buffer")
		//return false
	}
	t.OutQ.Push(flu)
	flu.ConfirmReceive()
	plog.Info("transformation out", flu.ID)
	flu_logger_svc.LogStepExit(flu.FeedLineUnit, step_type.Transformation, flu.Redelivered())

	return true
}

func newStdTransformer() *transformationStep {
	ts := &transformationStep{
		Step:          step.New(step_type.Transformation),
		stepConfigSvc: work_flow_svc.NewStepConfigService(),
	}

	ts.SetFluProcessor(ts.processFlu)
	return ts
}
