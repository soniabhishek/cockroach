package status_codes

type StatusCode int

const (
	Success         StatusCode = 1000
	CallBackFailure StatusCode = 1001
	FluRespFailure  StatusCode = 1002
	UnknownFailure  StatusCode = 1111
)

const (
	FF_FluIdNotPresent = "FF_0001"
	FF_RefIdNotPresent = "FF_0002"
	FF_TagIdNotPresent = "FF_0003"
	FF_ResultInvalid   = "FF_0004"
	FF_Other           = "FF_0000"
)
