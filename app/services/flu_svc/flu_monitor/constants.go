package flu_monitor

const (
	//TODO qc-uuid configure in db for active projects
	//TODO change this to file config instead of db
	HMAC_HEADER_KEY = "hmac_header"
	HMAC_KEY        = "hmac_key"
	IS_HMAC_ENABLED = "is_hmac_enabled"

	MAX_FLU_COUNT = "max_flu_count"
	CLIENT_QPS    = "client_qps"

	STATUS_COMPLETED = "COMPLETED"
	RESULT           = "result"
)
