package validation_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/clients"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
	"github.com/crowdflux/angel/app/services/work_flow_executor_svc/step"
	"github.com/crowdflux/angel/app/services/work_flow_svc"
)

type validationStep struct {
	step.Step
	stepConfigSvc work_flow_svc.IStepConfigSvc
}

func (v *validationStep) processFlu(flu feed_line.FLU) {
	v.AddToBuffer(flu)
	plog.Info("validation Step flu reached", flu.ID)

	vStep, err := v.stepConfigSvc.GetValidationStepConfig(flu.StepId)
	if err != nil {
		plog.Error("validation step", err, "fluId: "+flu.ID.String(), "stepid: "+flu.StepId.String(), flu.FeedLineUnit)
		flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.Validation, "ValidationConfigError", flu.Redelivered())
		return
	}

	isValidated, err := clients.GetMegatronClient().Validate(flu.Build, vStep.TemplateId)
	if err != nil {
		plog.Error("validation step", err, "fluId: "+flu.ID.String(), flu.FeedLineUnit)
		flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.Validation, "ValidationError", flu.Redelivered())
		return
	}

	flu.Build[vStep.AnswerKey] = isValidated

	v.finishFlu(flu)

}

func (v *validationStep) finishFlu(flu feed_line.FLU) bool {

	err := v.RemoveFromBuffer(flu)
	if err != nil {
		plog.Trace("validation step", "flu not present in buffer")
		//return false
	}
	v.OutQ.Push(flu)
	flu.ConfirmReceive()
	plog.Info("validation out", flu.ID)
	flu_logger_svc.LogStepExit(flu.FeedLineUnit, step_type.Validation, flu.Redelivered())

	return true
}

func newStdValidator() *validationStep {
	ts := &validationStep{
		Step:          step.New(step_type.Validation),
		stepConfigSvc: work_flow_svc.NewStepConfigService(),
	}

	ts.SetFluProcessor(ts.processFlu)
	return ts
}
