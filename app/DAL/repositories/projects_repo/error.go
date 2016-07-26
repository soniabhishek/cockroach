package projects_repo

import "github.com/crowdflux/angel/app/services/plerrors"

var ErrProjectNotFound = plerrors.ServiceError{"PR_0001", "Project not found"}
