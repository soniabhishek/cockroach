package flu_svc

import "gitlab.com/playment-main/angel/app/services/plerrors"

var ErrDuplicateReferenceId plerrors.ServiceError = plerrors.ServiceError{"FS_0001", "Duplicate Reference Id"}
var ErrReferenceIdMissing plerrors.ServiceError = plerrors.ServiceError{"FS_0002", "Reference Id Missing"}
var ErrDataValidation plerrors.ServiceError = plerrors.ServiceError{"FS_0003", "Invalid data passed"}
var ErrDataMissing plerrors.ServiceError = plerrors.ServiceError{"FS_0004", "Data missing"}
var ErrTagMissing plerrors.ServiceError = plerrors.ServiceError{"FS_0005", "Tag Missing"}
var ErrFluNotFound plerrors.ServiceError = plerrors.ServiceError{"FS_0006", "FeedLineUnit not found"}

type DataValidationError struct {
	plerrors.ServiceError
	Validations []validationError `json:"validations"`
}
