package acl_type

type ACL_Request struct {
	Token string   `json:"token"`
	Roles []string `json:"roles"`
}

type ACL_Response struct {
	Success  bool                 `json:"success"`
	Response ACL_Valid_Response   `json:"response"`
	Error    ACL_Failure_Response `json:"error"`
}
type ACL_Valid_Response struct {
	IsPermitted bool `json:"is_permitted"`
}

type ACL_Failure_Response struct {
	ErrorCode bool   `json:"error_code"`
	Message   string `json:"message"`
}
