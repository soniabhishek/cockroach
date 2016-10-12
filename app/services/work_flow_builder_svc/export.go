package work_flow_builder_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/step_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/step_router_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/workflow_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IWorkflowBuilderService interface {
	GetWorkflowContainer(uuid.UUID) (models.WorkflowContainer, error)
	AddWorkflowContainer(models.WorkflowContainer) (models.WorkflowContainer, error)
	InitWorkflowContainer(uuid.UUID) (models.WorkflowContainer, error)
	UpdateWorkflowContainer(models.WorkflowContainer) (models.WorkflowContainer, error)
}

func New() IWorkflowBuilderService {
	return &workFlowBuilderService{
		stepRepo:       step_repo.New(),
		stepRouterRepo: step_router_repo.New(),
		workflowRepo:   workflow_repo.New(),
		projectsRep:    projects_repo.New(),
	}
}
