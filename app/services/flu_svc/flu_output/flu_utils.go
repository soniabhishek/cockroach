package flu_output

import (
	"encoding/json"
	"fmt"
	"gitlab.com/playment-main/angel/app/models/status_codes"
	"io/ioutil"
	"net/http"
)

func ParseFluResponse(resp *http.Response) *Response {
	fluResp := &Response{}
	fluResp.HttpStatusCode = resp.StatusCode

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Println("response Headers:", resp)
	fmt.Println("response Body:", string(body))
	err := json.Unmarshal(body, fluResp)
	if err != nil {
		//TODO what to do with error
		fmt.Println(err)
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
