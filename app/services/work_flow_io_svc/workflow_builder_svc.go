package work_flow_builder_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/step_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/step_router_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/workflow_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/lib/pq"
	"time"
)

type workFlowBuilderService struct {
	stepRepo       step_repo.IStepRepo
	stepRouterRepo step_router_repo.IStepRoutesRepo
	workflowRepo   workflow_repo.IWorkflowRepo
	projectsRep    projects_repo.IProjectsRepo
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
	return models.WorkflowContainer{
		workflow,
		steps,
		routes,
	}, nil

}

func (w *workFlowBuilderService) InitWorkflowContainer(projectId uuid.UUID) (workflowContainer models.WorkflowContainer, err error) {
	exist, err := w.projectsRep.IfIdExist(projectId)
	if err != nil {
		return
	}
	if !exist {
		err = projects_repo.ErrProjectNotFound
		return
	}
	newWorkflow := models.WorkFlow{}
	newWorkflow.ID = uuid.NewV4()
	newWorkflow.CreatedAt = pq.NullTime{time.Now(), true}
	newWorkflow.UpdatedAt = newWorkflow.CreatedAt
	newWorkflow.ProjectId = projectId

	cornerSteps := generateCornerSteps(newWorkflow.ID)

	err = w.workflowRepo.Add(&newWorkflow)
	if err != nil {
		return
	}
	err = w.stepRepo.AddMany(cornerSteps)
	if err != nil {
		return
	}
	workflowContainer.WorkFlow = newWorkflow
	workflowContainer.Steps = cornerSteps
	return

}

func (w *workFlowBuilderService) AddWorkflowContainer(receivedWorkflowContainer models.WorkflowContainer) (workflowContainer models.WorkflowContainer, err error) {
	exist, err := w.projectsRep.IfIdExist(receivedWorkflowContainer.ProjectId)
	if err != nil {
		return
	}
	if !exist {
		err = projects_repo.ErrProjectNotFound
		return
	}
	receivedWorkflowContainer.WorkFlow.ID = uuid.NewV4()
	receivedWorkflowContainer.WorkFlow.CreatedAt = pq.NullTime{time.Now(), true}
	receivedWorkflowContainer.WorkFlow.UpdatedAt = receivedWorkflowContainer.WorkFlow.CreatedAt

	creationTime := pq.NullTime{time.Now(), true}
	for i, _ := range receivedWorkflowContainer.Steps {
		receivedWorkflowContainer.Steps[i].WorkFlowId = receivedWorkflowContainer.WorkFlow.ID
		receivedWorkflowContainer.Steps[i].CreatedAt = creationTime
		receivedWorkflowContainer.Steps[i].UpdatedAt = creationTime
	}
	for i, _ := range receivedWorkflowContainer.Routes {
		receivedWorkflowContainer.Routes[i].CreatedAt = creationTime
		receivedWorkflowContainer.Routes[i].UpdatedAt = creationTime
	}
	/* Here All Validations Needs to be Performed so as to make Complete Operation Atomic*/

	err = w.workflowRepo.Add(&receivedWorkflowContainer.WorkFlow) /*Perform in End */
	if err != nil {
		return
	}
	err = w.stepRepo.AddMany(receivedWorkflowContainer.Steps)
	if err != nil {
		return
	}
	err = w.stepRouterRepo.AddMany(receivedWorkflowContainer.Routes)
	if err != nil {
		return
	}
	return receivedWorkflowContainer, nil
}

func (w *workFlowBuilderService) UpdateWorkflowContainer(receivedWorkflowContainer models.WorkflowContainer) (workflowContainer models.WorkflowContainer, err error) {
	exist, err := w.projectsRep.IfIdExist(receivedWorkflowContainer.ProjectId)
	if err != nil {
		return
	}
	if !exist {
		err = projects_repo.ErrProjectNotFound
		return
	}
	exist, err = w.workflowRepo.IfIdExist(receivedWorkflowContainer.ID)
	if err != nil || !exist {
		err = workflow_repo.ErrWorkflowNotFound
		return
	}
	existingSteps, err := w.stepRepo.GetStepsByWorkflowId(receivedWorkflowContainer.ID)
	if err != nil {
		return
	}
	existingRoutes, err := w.stepRouterRepo.GetRoutesByWorkFlowId(receivedWorkflowContainer.ID)
	if err != nil {
		return
	}

	insertRoutes, updateRoutes, deleteRoutes, err := computeRouteComparision(receivedWorkflowContainer.Routes, existingRoutes)
	if err != nil {
		return
	}
	insertSteps, updateSteps, deleteSteps, err := computeStepComparision(receivedWorkflowContainer.Steps, existingSteps, receivedWorkflowContainer.ID)
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

	return w.GetWorkflowContainer(receivedWorkflowContainer.ID)

}

func computeStepComparision(receivedSteps, existingSteps []models.Step, workflowID uuid.UUID) ([]models.Step, []models.Step, []models.Step, error) {
	var forUpdate, forInsert []models.Step
	for _, received := range receivedSteps {
		var update bool
		for index, existing := range existingSteps {
			if received.ID == existing.ID {
				update = true
				//compare received and existing value here
				received.UpdatedAt = pq.NullTime{time.Now(), true}
				received.WorkFlowId = workflowID
				forUpdate = append(forUpdate, received)
				existingSteps = append(existingSteps[:index], existingSteps[index+1:]...)
				break
			}
		}
		if !update {
			creationTime := pq.NullTime{time.Now(), true}
			received.CreatedAt = creationTime
			received.UpdatedAt = creationTime
			received.WorkFlowId = workflowID
			forInsert = append(forInsert, received)
		}
	}
	return forInsert, forUpdate, existingSteps, nil
}

func computeRouteComparision(receivedRoutes, existingRoutes []models.Route) ([]models.Route, []models.Route, []models.Route, error) {
	var forUpdate, forInsert []models.Route
	for _, received := range receivedRoutes {
		var update bool
		for index, existing := range existingRoutes {
			if received.ID == existing.ID {
				update = true
				received.UpdatedAt = pq.NullTime{time.Now(), true}
				forUpdate = append(forUpdate, received)
				existingRoutes = append(existingRoutes[:index], existingRoutes[index+1:]...)
				break
			}
		}
		if !update {
			creationTime := pq.NullTime{time.Now(), true}
			received.CreatedAt = creationTime
			received.UpdatedAt = creationTime
			forInsert = append(forInsert, received)
		}
	}
	return forInsert, forUpdate, existingRoutes, nil
}
func generateCornerSteps(workflowId uuid.UUID) []models.Step {
	startStep := models.Step{}
	startStep.ID = uuid.NewV4()
	startStep.CreatedAt = pq.NullTime{time.Now(), true}
	startStep.UpdatedAt = startStep.CreatedAt
	startStep.IsStart = true
	startStep.Type = step_type.Manual
	startStep.WorkFlowId = workflowId

	endStep := models.Step{}
	endStep.ID = uuid.NewV4()
	endStep.CreatedAt = pq.NullTime{time.Now(), true}
	endStep.UpdatedAt = endStep.CreatedAt
	endStep.Type = step_type.Manual
	endStep.WorkFlowId = workflowId

	return []models.Step{startStep, endStep}
}
