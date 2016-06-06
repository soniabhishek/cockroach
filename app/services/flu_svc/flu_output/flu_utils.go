package flu_output

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.com/playment-main/angel/app/models/status_codes"
	"gitlab.com/playment-main/angel/app/plog"
)

func ParseFluResponse(resp *http.Response) *Response {
	fluResp := &Response{}
	fluResp.HttpStatusCode = resp.StatusCode

	body, _ := ioutil.ReadAll(resp.Body)
	plog.Info("response Status:", resp.Status)
	plog.Info("response Headers:", resp.Header)
	plog.Info("response Headers:", resp)
	plog.Info("response Body:", string(body))
	err := json.Unmarshal(body, fluResp)
	if err != nil {
		plog.Error("Response Parsing Error: ", err)
		return fluResp
	}
	return fluResp
}

/* This function will check whether we need to try calling again to the hitting server.
It calls back in case of
-	HTTP ERR 408 (Request Timeout)
-	HTTP ERR 504 (Gateway Timeout)
*/
func HttpCodeForCallback(httpStatusCode int) bool {
	switch httpStatusCode {
	case
		http.StatusRequestTimeout,
		http.StatusGatewayTimeout:
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
