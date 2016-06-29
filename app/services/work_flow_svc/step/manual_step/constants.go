package manual_step

const (
	ID         = "Id"
	REF_ID     = "ReferenceId"
	DATA       = "Data"
	BUILD      = "Build"
	TAG        = "Tag"
	PROJECT_ID = "ProjectId"
	STEP_ID    = "StepId"
	CREATED_ID = "CreatedAt"
	UPDATED_AT = "UpdatedAt"
)

const (
	//TODO find an alternative. Might be like make it configurable.
	TEMP_FOLDER  = `./`
	FLU_ID_INDEX = 0
	BUILD_INDEX  = 1
	COMMA        = ","
)
const (
	DOWNLOAD_ENDPOINT = `/manual_step_download`
	UPLOAD_ENDPOINT   = `/manual_step_upload`

	MANUAL_STEP_ID    = "manualStepId"
	UPLOAD            = "upload"
	PARAM_FILES       = "files"
	PARAM_PLAYMENT_ID = "playment_request_id"

	CONTENT_TYPE = "Content-Type"
	SUCCESS      = "success"
	ERROR        = "error"
	FILEPATH     = "filepath"
	URL          = "url"
)
