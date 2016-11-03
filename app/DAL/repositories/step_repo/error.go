package step_repo

import "github.com/crowdflux/angel/app/services/plerrors"

var ErrSteptNotFound = plerrors.ServiceError{"SR_0001", "Step not found"}
