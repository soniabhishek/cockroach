package flu_validator

import "github.com/crowdflux/angel/app/services/plerrors"

var errDataValidation plerrors.ServiceError = plerrors.ServiceError{"FS_0003", "Invalid data passed"}

type DataValidationError struct {
	plerrors.ServiceError
	Validations []validationError `json:"validations"`
}
