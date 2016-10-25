package work_flow_io_svc

import (
	"fmt"
	"github.com/crowdflux/angel/app/DAL/repositories/step_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/step_router_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/workflow_repo"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWorkFlowBuilderService_GetWorkflowContainer(t *testing.T) {

	serviceObject := workFlowBuilderService{}
	serviceObject.workflowRepo = workflow_repo.New()
	serviceObject.stepRouterRepo = step_router_repo.New()
	serviceObject.stepRepo = step_repo.New()

	container, err := serviceObject.GetWorkflowContainer(uuid.NewV4())
	fmt.Println(container, err)

	assert.NoError(t, err)
}

/**
Checks and test to be written

1. If in update method both projectId and WorkflowId are valid but not associated with each other based on foreign key constraint
*/
