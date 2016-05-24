package flu_validator

import "gitlab.com/playment-main/angel/app/services/plerrors"

var ErrDataValidation plerrors.ServiceError = plerrors.ServiceError{"FS_0003", "Invalid data passed"}

type DataValidationError struct {
	plerrors.ServiceError
	Validations []validationError `json:"validations"`
}
