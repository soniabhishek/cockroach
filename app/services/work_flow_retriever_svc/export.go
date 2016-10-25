package work_flow_retriever_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/workflow_repo"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/models"
)

type IWorkflowRetrieverService interface {
	GetWorkFlow(projectId uuid.UUID, tag string) ([]models.WorkFlow, error)
}

func New() IWorkflowRetrieverService {
	return &workFlowRetrieverService{
	 	workflowRepo:   workflow_repo.New(),
	}
}


