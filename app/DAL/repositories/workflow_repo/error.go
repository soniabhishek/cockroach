package workflow_repo

import "github.com/crowdflux/angel/app/services/plerrors"

var ErrWorkflowNotFound = plerrors.ServiceError{"WR_0001", "WorkFlow not found"}
