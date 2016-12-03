package work_flow_io_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/clients_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/step_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/step_router_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/workflow_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/workflow_tags_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IWorkflowBuilderService interface {
	GetWorkflowContainer(uuid.UUID) (models.WorkflowContainer, error)
	AddWorkflowContainer(models.WorkflowContainer) (models.WorkflowContainer, error)
	UpdateWorkflowContainer(models.WorkflowContainer) (models.WorkflowContainer, error)
	CloneWorkflowContainer(models.WorkFlowCloneModel) (models.WorkflowContainer, error)
	FetchWorkflowsByProjectId(uuid.UUID) ([]models.WorkFlow, error)
}

func New() IWorkflowBuilderService {
	return &workFlowBuilderService{
		stepRepo:         step_repo.New(),
		stepRouterRepo:   step_router_repo.New(),
		workflowRepo:     workflow_repo.New(),
		projectsRep:      projects_repo.New(),
		workflowTagsRepo: workflow_tags_repo.New(),
		clientsRepo:      clients_repo.New(),
	}
}
