package transformation_step

import (
	"gitlab.com/playment-main/angel/app/DAL/repositories/projects_repo"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
)

func newStdTransformer() *transformationStep {
	return &transformationStep{
		Step:         step.New(),
		projectsRepo: projects_repo.New(),
	}
}

var StdTransformationStep = newStdTransformer()
