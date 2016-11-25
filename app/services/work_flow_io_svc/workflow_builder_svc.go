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
	"github.com/pkg/errors"
)

type workFlowBuilderService struct {
	stepRepo         step_repo.IStepRepo
	stepRouterRepo   step_router_repo.IStepRoutesRepo
	workflowRepo     workflow_repo.IWorkflowRepo
	projectsRep      projects_repo.IProjectsRepo
	workflowTagsRepo workflow_tags_repo.IWorkflowTagsRepo
	clientsRepo      clients_repo.IClientsRepo
}

var _ IWorkflowBuilderService = &workFlowBuilderService{}

func (w *workFlowBuilderService) GetWorkflowContainer(workflowId uuid.UUID) (workflowContainer models.WorkflowContainer, err error) {
	workflow, err := w.workflowRepo.GetById(workflowId)
	if err != nil {
		return
	}
	steps, err := w.stepRepo.GetStepsByWorkflowId(workflowId)
	if err != nil {
		return
	}
	routes, err := w.stepRouterRepo.GetRoutesByWorkFlowId(workflowId)
	if err != nil {
		return
	}
	tags, err := w.workflowTagsRepo.GetByWorkFlowId(workflowId)
	if err != nil {
		return
	}
	return models.WorkflowContainer{
		workflow,
		steps,
		routes,
		tags,
	}, nil

}

/**
This will be used to add a new Workflow for a given projectId

*/
func (w *workFlowBuilderService) AddWorkflowContainer(receivedWorkflowContainer models.WorkflowContainer) (workflowContainer models.WorkflowContainer, err error) {
	//Checks If the project Id is valid or not
	exist, err := w.projectsRep.IfIdExist(receivedWorkflowContainer.ProjectId)
	if err != nil {
		return
	}
	if !exist {
		err = projects_repo.ErrProjectNotFound
		return
	}
	if len(receivedWorkflowContainer.Tags) < 1 {
		err = errors.New("Tags list is empty")
		return
	}

	/* Here All Validations Needs to be Performed so as to make Complete Operation Atomic*/

	err = w.workflowRepo.Add(&receivedWorkflowContainer.WorkFlow) /*Perform in End */
	if err != nil {
		return
	}

	//Here we will update the creation and updated at time

	for i, _ := range receivedWorkflowContainer.Tags {
		receivedWorkflowContainer.Tags[i].WorkFlowId = receivedWorkflowContainer.WorkFlow.ID
		receivedWorkflowContainer.Tags[i].ProjectId = receivedWorkflowContainer.ProjectId
	}

	for i, _ := range receivedWorkflowContainer.Steps {
		receivedWorkflowContainer.Steps[i].WorkFlowId = receivedWorkflowContainer.WorkFlow.ID
	}

	err = w.stepRepo.AddMany(receivedWorkflowContainer.Steps)
	if err != nil {
		return
	}
	err = w.stepRouterRepo.AddMany(receivedWorkflowContainer.Routes)
	if err != nil {
		return
	}
	err = w.workflowTagsRepo.Add(receivedWorkflowContainer.Tags)
	if err != nil {
		return
	}
	return receivedWorkflowContainer, nil
}

/**
This will update the existing workflow
*/
func (w *workFlowBuilderService) UpdateWorkflowContainer(receivedWorkflowContainer models.WorkflowContainer) (workflowContainer models.WorkflowContainer, err error) {
	//This will check if Project Id in body exist or not
	exist, err := w.projectsRep.IfIdExist(receivedWorkflowContainer.ProjectId)
	if err != nil {
		return
	}
	if !exist {
		err = projects_repo.ErrProjectNotFound
		return
	}
	if len(receivedWorkflowContainer.Tags) < 1 {
		err = errors.New("Tags list is empty")
		return
	}
	//This will check if Workflow Id in body exist or not
	exist, err = w.workflowRepo.IfIdExist(receivedWorkflowContainer.ID)
	if err != nil {
		return
	}
	if !exist {
		err = workflow_repo.ErrWorkflowNotFound
		return
	}

	//existing steps will be fetched based on the provided workflowId
	existingSteps, err := w.stepRepo.GetStepsByWorkflowId(receivedWorkflowContainer.ID)
	if err != nil {
		return
	}
	//existing routes will be fetched based on the provided workflowId
	existingRoutes, err := w.stepRouterRepo.GetRoutesByWorkFlowId(receivedWorkflowContainer.ID)
	if err != nil {
		return
	}
	//existing tags will be fetched based on the provided workflowId
	existingTags, err := w.workflowTagsRepo.GetByWorkFlowId(receivedWorkflowContainer.ID)
	if err != nil {
		return
	}
	//This will categorize by comparing existing routes with routes in body into what needs to be inserted, updates or deleted
	insertRoutes, updateRoutes, deleteRoutes, err := computeRouteComparision(receivedWorkflowContainer.Routes, existingRoutes)
	if err != nil {
		return
	}

	//This will categorize by comparing existing steps with steps in body into what needs to be inserted, updates or deleted
	insertSteps, updateSteps, deleteSteps, err := computeStepComparision(receivedWorkflowContainer.Steps, existingSteps, receivedWorkflowContainer.ID)
	if err != nil {
		return
	}

	//This will categorize by comparing existing steps with steps in body into what needs to be inserted, updates or deleted
	insertTags, updateTags, deleteTags, err := computeTagsComparision(receivedWorkflowContainer.Tags, existingTags, receivedWorkflowContainer.ID, receivedWorkflowContainer.ProjectId)
	if err != nil {
		return
	}
	//Sequence should be add, update, delete

	err = w.stepRepo.AddMany(insertSteps)
	if err != nil {
		return
	}

	err = w.stepRouterRepo.AddMany(insertRoutes)
	if err != nil {
		return
	}

	_, err = w.stepRepo.UpdateMany(updateSteps)
	if err != nil {
		return
	}

	_, err = w.stepRouterRepo.UpdateMany(updateRoutes)
	if err != nil {
		return
	}

	_, err = w.stepRouterRepo.DeleteMany(deleteRoutes)
	if err != nil {
		return
	}

	_, err = w.stepRepo.DeleteMany(deleteSteps)
	if err != nil {
		return
	}

	err = w.workflowTagsRepo.Add(insertTags)
	if err != nil {
		return
	}

	err = w.workflowTagsRepo.Update(updateTags)
	if err != nil {
		return
	}

	err = w.workflowTagsRepo.Delete(deleteTags)
	if err != nil {
		return
	}
	err = w.workflowRepo.Update(&receivedWorkflowContainer.WorkFlow)
	if err != nil {
		return
	}
	//finally after all insert update and delete mechanism we will fetch whole new workflow from backend
	return w.GetWorkflowContainer(receivedWorkflowContainer.ID)

}

/**
This will clone a workflow
*/
func (w *workFlowBuilderService) CloneWorkflowContainer(workflowCloneData models.WorkFlowCloneModel) (clonedContainer models.WorkflowContainer, err error) {
	//Checks If the project Id is valid or not

	exist, err := w.workflowRepo.IfIdExist(workflowCloneData.WorkFlowId)
	if err != nil {
		return
	}
	if !exist {
		err = workflow_repo.ErrWorkflowNotFound
		return
	}
	exist, err = w.clientsRepo.IfIdExist(workflowCloneData.ClientId)
	if err != nil {
		return
	}
	if !exist {
		err = workflow_repo.ErrWorkflowNotFound
		return
	}
	existingContainer, err := w.GetWorkflowContainer(workflowCloneData.WorkFlowId)
	if err != nil {
		return
	}
	existingContainer.WorkFlow.ProjectId = workflowCloneData.ProjectId
	existingContainer.WorkFlow.Label = workflowCloneData.Label
	existingContainer.Tags = workflowCloneData.Tags
	stepIdMapping := make(map[uuid.UUID]uuid.UUID)
	for i, step := range existingContainer.Steps {
		stepIdMapping[step.ID] = uuid.NewV4()
		existingContainer.Steps[i].ID = stepIdMapping[step.ID]
	}
	for i, route := range existingContainer.Routes {
		existingContainer.Routes[i].ID = uuid.NewV4()
		existingContainer.Routes[i].StepId = stepIdMapping[route.StepId]
		existingContainer.Routes[i].NextStepId = stepIdMapping[route.NextStepId]
	}
	clonedContainer, err = w.AddWorkflowContainer(existingContainer)
	if err != nil {
		return
	}
	return
}
