package step_router_repo

import "gitlab.com/playment-main/angel/app/services/plerrors"

var ErrRouteNotFound = plerrors.ServiceError{"SR_0001", "Route not found"}
