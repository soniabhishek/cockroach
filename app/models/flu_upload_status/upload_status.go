package flu_upload_status

type FluUploadStatus uint

const (
	Pending FluUploadStatus = iota
	Processing
	Success
	Failure
	PartialUpload
)
