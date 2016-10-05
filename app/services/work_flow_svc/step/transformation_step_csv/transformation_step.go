package transformation_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/clients"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/step_configuration_repo"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
)

type transformationStep struct {
	step.Step
	transformationConfigRepo step_configuration_repo.ITransformationStepConfigurationRepo
}

func (t *transformationStep) processFlu(flu feed_line.FLU) {
	t.AddToBuffer(flu)
	plog.Info("transformation Step flu reached", flu)

	tStep, err := t.transformationConfigRepo.GetByStepId(flu.StepId)
	if err != nil {
		plog.Error("transformation step", err)
		return
	}

	transformedBuild, err := clients.GetMegatronClient().Transform(flu.Build, tStep.TemplateId)
	if err != nil {
		plog.Error("Transformation step", err)
		return
	}

	plog.Info("transformation step", transformedBuild)

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
	plog.Info("transformation out", flu.ID)
	return true
}

func newStdTransformer() *transformationStep {
	ts := &transformationStep{
		Step: step.New(step_type.Transformation),
		transformationConfigRepo: step_configuration_repo.NewTransformationStepConfigurationRepo(),
	}

	ts.SetFluProcessor(ts.processFlu)
	return ts
}
