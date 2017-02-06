package flu_monitor

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/crowdflux/angel/app/models/status_codes"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/plog/log_tags"
)

type WebhookResponse struct {
	HttpStatusCode int
	//FluStatusCode  status_codes.StatusCode
	Invalid_Flus []invalidFlu `json:"invalid_flus"`
	RawResponse  string
}

func ParseFluResponse(resp *http.Response) *WebhookResponse {
	fluResp := &WebhookResponse{}
	fluResp.HttpStatusCode = resp.StatusCode

	body, _ := ioutil.ReadAll(resp.Body)
	fluResp.RawResponse = string(body)

	err := json.Unmarshal(body, fluResp)
	if err != nil {

		/*TODO How to handle this error.
		errors.New("asfdsgdf")

		switch err.(type) {
		case :

		}*/

		plog.Error("Response Parsing Error: ", err, plog.MP(log_tags.POSTBACK_RESPONSE, fluResp))
		return fluResp

	}
	return fluResp
}

/* This function will check whether we need to try calling again to the hitting server.
It calls back in case of
-	HTTP ERR 408 (Request Timeout)
-	HTTP ERR 500 (Internal Server Error)
-	HTTP ERR 501 (Not Implemented)
-	HTTP ERR 502 (Bad Gateway)
-	HTTP ERR 503 (Service Unavailable)
-	HTTP ERR 504 (Gateway Timeout)
-	HTTP ERR 505 (HTTP Version Not Supported)
-	HTTP ERR 511 (Network Authentication Required)
*/
func HttpCodeForCallback(httpStatusCode int) bool {
	switch httpStatusCode {
	case
		http.StatusRequestTimeout,
		http.StatusInternalServerError,
		http.StatusNotImplemented,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusHTTPVersionNotSupported,
		http.StatusNetworkAuthenticationRequired:
		return true
	}
	return false
}

func IsValidInternalError(internalCode string) bool {
	switch internalCode {
	case
		status_codes.FF_FluIdNotPresent,
		status_codes.FF_RefIdNotPresent,
		status_codes.FF_TagIdNotPresent,
		status_codes.FF_ResultInvalid,
		status_codes.FF_Other:
		return true
	}
	return false
}
