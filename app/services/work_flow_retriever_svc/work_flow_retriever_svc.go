package work_flow_retriever_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/workflow_repo"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/models"
)

type workFlowRetrieverService struct {
	workflowRepo   workflow_repo.IWorkflowRepo
}

func (w *workFlowRetrieverService) GetWorkFlow(projectId uuid.UUID, tag string) (workflow []models.WorkFlow, err error){

	if tag=="" {
		workflow, err = w.workflowRepo.GetWorkFlowsByProjectId(projectId)
		if err != nil {
			return
		}
	}else{
		wf, err := w.workflowRepo.GetWorkFlowByProjectIdAndTag(projectId, tag)
		if err != nil {
			return workflow, err
		}
		workflow = append(workflow, wf)
	}
	return
}