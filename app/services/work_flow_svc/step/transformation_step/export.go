package transformation_step

import (
	"github.com/crowdflux/angel/app/DAL/repositories/step_configuration_repo"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
)

func newStdTransformer() *transformationStep {
	return &transformationStep{
		Step: step.New(),
		transformationConfigRepo: step_configuration_repo.NewTransformationStepConfigurationRepo(),
	}
}

var StdTransformationStep = newStdTransformer()
