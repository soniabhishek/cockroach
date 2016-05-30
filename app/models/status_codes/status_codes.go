package status_codes

type StatusCode int
const (
	Success = 1000
	CallBackFailure = 1001
	FluRespFailure = 1002
	UnknownFailure = 1111
)

const (

	FF_FluIdNotPresent = "FF_0001"
	FF_RefIdNotPresent = "FF_0002"
	FF_TagIdNotPresent = "FF_0003"
	FF_ResultInvalid = "FF_0004"
	FF_Other = "FF_0000"
)