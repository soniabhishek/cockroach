package flu_monitor

import (
	"github.com/crowdflux/angel/app/models/status_codes"
	"github.com/crowdflux/angel/app/plog"
	"net/http"
)

func validationErrorCallback(resp *http.Response) (*FluResponse, status_codes.StatusCode) {
	defer resp.Body.Close()

	fluResp := ParseFluResponse(resp)
	shouldCallBack := HttpCodeForCallback(fluResp.HttpStatusCode)
	plog.Trace("HTTPStatusCode: [", fluResp.HttpStatusCode, "] Should Call back: ", shouldCallBack)
	if shouldCallBack {
		return fluResp, status_codes.CallBackFailure
	} else {
		//If any invalid flu response code is in our InvalidationCodeArray, then we log[ERROR] it
		for _, invalidFlu := range fluResp.Invalid_Flus {
			if IsValidInternalError(invalidFlu.Error) {
				return fluResp, status_codes.FluRespFailure
			}
		}
	}
	return fluResp, status_codes.Success
}
