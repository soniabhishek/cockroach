package transformation_step

import (
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
)

func newStdTransformer() *transformationStep {
	return &transformationStep{
		Step:         step.New(),
		projectsRepo: projects_repo.New(),
	}
}

var StdTransformationStep = newStdTransformer()
