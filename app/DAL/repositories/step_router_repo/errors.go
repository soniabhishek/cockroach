package step_router_repo

import "github.com/crowdflux/angel/app/services/plerrors"

var ErrRouteNotFound = plerrors.ServiceError{"SR_0001", "Route not found"}
