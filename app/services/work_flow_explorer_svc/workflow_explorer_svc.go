package work_flow_explorer_svc

import (
	"database/sql"
	"github.com/crowdflux/angel/app/DAL/repositories/clients_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/user_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/workflow_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/utilities/clients/validator"
	"github.com/pkg/errors"
)

type workflowExplorerService struct {
	clientsRepo  clients_repo.IClientsRepo
	userRepo     user_repo.IUserRepo
	projectsRepo projects_repo.IProjectsRepo
	workflowRepo workflow_repo.IWorkflowRepo
}

var _ IWorkFlowExplorerService = &workflowExplorerService{}

func (ws *workflowExplorerService) CreateClient(client models.Client) (response models.Client, err error) {
	if client.Name == "" {
		err = errors.New("Invalid Client Name")
		return
	}
	user := models.User{}
	user.Username = client.Name
	user.FirstName = sql.NullString{client.Name, true}

	err = ws.userRepo.Add(&user)
	if err != nil {
		return
	}
	client.UserId = user.ID
	err = ws.clientsRepo.Add(&client)
	if err != nil {
		return
	}
	client.ClientSecretUuid = uuid.Nil
	return client, nil
}

func (ws *workflowExplorerService) GetClient(uuid.UUID) (client models.Client, err error) {
	return
}

func (ws *workflowExplorerService) FetchAllClient() ([]models.Client, error) {
	return ws.clientsRepo.GetAllClients()
}

func (ws *workflowExplorerService) FetchProjectsByClientId(clientId uuid.UUID) ([]models.Project, error) {
	return ws.projectsRepo.GetByClientId(clientId)
}

func (ws *workflowExplorerService) FetchWorkflowsByProjectId(projectId uuid.UUID) ([]models.WorkFlow, error) {
	return ws.workflowRepo.GetWorkFlowsByProjectId(projectId)
}

func (ws *workflowExplorerService) CreateProject(project models.Project) (response models.Project, err error) {
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
