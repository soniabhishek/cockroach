package work_flow_explorer_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/clients_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/user_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/workflow_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/utilities/clients/validator"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"time"
)

type workflowExplorerService struct {
	clientsRepo  clients_repo.IClientsRepo
	userRepo     user_repo.IUserRepo
	projectsRepo projects_repo.IProjectsRepo
	workflowRepo workflow_repo.IWorkflowRepo
}

func (cs *workflowExplorerService) CreateClient(client models.Client) (response models.Client, err error) {
	err = validator.ValidateClient(client)
	if err != nil {
		return
	}
	exist, err := cs.userRepo.IfIdExist(client.UserId)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New("Invalid User Id")
		return
	}
	client.ClientSecretUuid = uuid.NewV4()
	client.ID = uuid.NewV4()
	currentTime := time.Now()
	client.CreatedAt = pq.NullTime{currentTime, true}
	client.UpdatedAt = pq.NullTime{currentTime, true}
	err = cs.clientsRepo.Add(client)
	if err != nil {
		return
	}
	return client, nil
}

func (cs *workflowExplorerService) GetClient(uuid.UUID) (client models.Client, err error) {
	return
}

func (cs *workflowExplorerService) FetchAllClient() ([]models.ClientModel, error) {
	return cs.clientsRepo.GetAllClients()
}

func (cs *workflowExplorerService) FetchProjectsByClientId(clientId uuid.UUID) ([]models.Project, error) {
	return cs.projectsRepo.GetByClientId(clientId)
}

func (cs *workflowExplorerService) FetchWorkflowsByProjectId(projectId uuid.UUID) ([]models.WorkFlow, error) {
	return cs.workflowRepo.GetWorkFlowsByProjectId(projectId)
}
