package step_configuration_repo

import "github.com/crowdflux/angel/app/services/plerrors"

var ErrStepConfigNotFound = plerrors.ServiceError{"SCR_0001", "Configuration not found"}
