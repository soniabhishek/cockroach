package workflow_tags_repo

import "github.com/crowdflux/angel/app/services/plerrors"

var ErrWorkflowTagsNotFound = plerrors.ServiceError{"TG_0001", "Tage not found for Workflow"}
