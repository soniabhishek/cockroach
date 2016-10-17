package flu_output

import "github.com/crowdflux/angel/app/models/status_codes"

//TODO rest will be added later
type FluResponse struct {
	HttpStatusCode int
	FluStatusCode  status_codes.StatusCode
	Invalid_Flus   []invalidFlu `json:"invalid_flus"`
	RawResponse    string
}
type invalidFlu struct {
	Flu_Id  string `json:"flu_id"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
