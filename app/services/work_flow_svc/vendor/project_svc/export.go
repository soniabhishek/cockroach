package project_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/clients_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IWorkFlowProjetService interface {
	CreateProject(models.Project) (models.Project, error)
	FetchProjectsByClientId(uuid.UUID) ([]models.Project, error)
}

func New() IWorkFlowProjetService {
	return &workFlowProjectService{
		clientsRepo:  clients_repo.New(),
		projectsRepo: projects_repo.New(),
	}
}
