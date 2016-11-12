package work_flow_explorer_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/clients_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/user_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/workflow_tags_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"time"
)

type workflowExplorerService struct {
	clientsRepo      clients_repo.IClientsRepo
	userRepo         user_repo.IUserRepo
	projectsRepo     projects_repo.IProjectsRepo
	workflowTagsRepo workflow_tags_repo.IWorkflowTagsRepo
}

func (cs *workflowExplorerService) CreateClient(client models.Client) (response models.Client, err error) {
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

func (cs *workflowExplorerService) FetchWorkflowsByProjectId(projectId uuid.UUID) ([]models.TagExplorerModel, error) {
	wfta, err := cs.workflowTagsRepo.GetByProjectId(projectId)
	if err != nil {
		return nil, err
	}
	if len(wfta) < 1 {
		return nil, err
	}
	temp := make(map[uuid.UUID]string)
	for _, v := range wfta {
		if val, ok := temp[v.WorkFlowId]; ok {
			temp[v.WorkFlowId] = val + " " + v.TagName
		} else {
			temp[v.WorkFlowId] = v.TagName
		}
	}
	var result []models.TagExplorerModel
	for k, v := range temp {
		result = append(result, models.TagExplorerModel{v, k})
	}
	return result, err
}
