package macro_task_repo

import "gitlab.com/playment-main/angel/app/services/plerrors"

var ErrMacroTaskNotFound = plerrors.ServiceError{"MA_0001", "MacroTask not found"}
