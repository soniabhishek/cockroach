package work_flow_io_svc

import (
	"github.com/lib/pq"
	"time"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

/**
This is a utility function for comparing, validating and classifying Steps also it will change their Updated at time
This should be containing all the validation steps if required in future for steps
*/
func computeStepComparision(receivedSteps, existingSteps []models.Step, workflowID uuid.UUID) ([]models.Step, []models.Step, []models.Step, error) {
	var forUpdate, forInsert []models.Step
	for _, received := range receivedSteps {
		var update bool
		for index, existing := range existingSteps {
			if received.ID == existing.ID {
				update = true
				//compare received and existing value here
				received.WorkFlowId = workflowID
				forUpdate = append(forUpdate, received)
				existingSteps = append(existingSteps[:index], existingSteps[index+1:]...)
				break
			}
		}
		if !update {
			received.WorkFlowId = workflowID
			forInsert = append(forInsert, received)
		}
	}
	return forInsert, forUpdate, existingSteps, nil
}

/**
This is a utility function for comparing, validating and classifying Routes also it will change their Updated at time
This should be containing all the validation routes if required in future for routes
*/
func computeRouteComparision(receivedRoutes, existingRoutes []models.Route) ([]models.Route, []models.Route, []models.Route, error) {
	var forUpdate, forInsert []models.Route
	for _, received := range receivedRoutes {
		var update bool
		for index, existing := range existingRoutes {
			if received.ID == existing.ID {
				update = true
				forUpdate = append(forUpdate, received)
				existingRoutes = append(existingRoutes[:index], existingRoutes[index+1:]...)
				break
			}
		}
		if !update {
			forInsert = append(forInsert, received)
		}
	}
	return forInsert, forUpdate, existingRoutes, nil
}

/**
This is a utility function for comparing, validating and classifying Tags also it will change their Updated at time
This should be containing all the validation routes if required in future for routes
*/
func computeTagsComparision(receivedTags, existingTags []models.WorkFlowTagAssociator, workflowID, projectID uuid.UUID) ([]models.WorkFlowTagAssociator, []models.WorkFlowTagAssociator, []models.WorkFlowTagAssociator, error) {
	var forUpdate, forInsert []models.WorkFlowTagAssociator
	for _, received := range receivedTags {
		var update bool
		for index, existing := range existingTags {
			if received.TagName == existing.TagName {
				update = true
				received.CreatedAt = pq.NullTime{time.Now(), true}
				received.WorkFlowId = workflowID
				received.ProjectId = projectID
				forUpdate = append(forUpdate, received)
				existingTags = append(existingTags[:index], existingTags[index+1:]...)
				break
			}
		}
		if !update {
			creationTime := pq.NullTime{time.Now(), true}
			received.WorkFlowId = workflowID
			received.ProjectId = projectID
			received.CreatedAt = creationTime
			forInsert = append(forInsert, received)
		}
	}
	return forInsert, forUpdate, existingTags, nil
}
