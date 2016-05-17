package plerrors

type ServiceError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (s ServiceError) Error() string {
	return "Service Error : " + s.Code + " : " + s.Message
}

//--------------------------------------------------------------------------------//

type IncorrectUUIDError struct {
	ServiceError
	UUIDFields []string `json:"fields"`
}

func ErrIncorrectUUID(fields ...string) IncorrectUUIDError {
	return IncorrectUUIDError{
		ServiceError: ServiceError{"GE_0001", "Invalid UUID passed"},
		UUIDFields:   fields,
	}
}

type RequestParamMissingError struct {
	ServiceError
	UUIDFields []string `json:"fields"`
}

func ErrRequestParamMissing(fields ...string) RequestParamMissingError {
	return RequestParamMissingError{
		ServiceError: ServiceError{"GE_0002", "Parameter missing"},
		UUIDFields:   fields,
	}
}

//--------------------------------------------------------------------------------//

var ErrMalformedJson ServiceError = ServiceError{
	Code:    "GE_0001",
	Message: "Malformed json",
}
