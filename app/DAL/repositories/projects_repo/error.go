package projects_repo

import "gitlab.com/playment-main/angel/app/services/plerrors"

var ErrProjectNotFound = plerrors.ServiceError{"PR_0001", "Project not found"}
