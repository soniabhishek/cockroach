package work_flow_explorer_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/clients_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/user_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/workflow_tags_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IWorkFlowExplorerService interface {
	CreateClient(models.Client) (models.Client, error)
	GetClient(uuid.UUID) (models.Client, error)
	FetchAllClient() ([]models.ClientModel, error)
	FetchProjectsByClientId(uuid.UUID) ([]models.Project, error)
	FetchWorkflowsByProjectId(uuid.UUID) ([]models.TagExplorerModel, error)
}

func New() IWorkFlowExplorerService {
	return &workflowExplorerService{
		clientsRepo:      clients_repo.New(),
		userRepo:         user_repo.New(),
		projectsRepo:     projects_repo.New(),
		workflowTagsRepo: workflow_tags_repo.New(),
	}
}
