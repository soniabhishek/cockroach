package project_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/clients_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/utilities/clients/validator"
	"github.com/pkg/errors"
)

type workFlowProjectService struct {
	clientsRepo  clients_repo.IClientsRepo
	projectsRepo projects_repo.IProjectsRepo
}

func (ws *workFlowProjectService) FetchProjectsByClientId(clientId uuid.UUID) ([]models.Project, error) {
	return ws.projectsRepo.GetByClientId(clientId)
}

func (ws *workFlowProjectService) CreateProject(project models.Project) (response models.Project, err error) {
	err = validator.ValidateProject(project)
	if err != nil {
		return
	}
	exist, err := ws.clientsRepo.IfIdExist(project.ClientId)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New("Invalid Project")
		return
	}
	err = ws.projectsRepo.Add(&project)
	if err != nil {
		return
	}
	return project, nil
}
