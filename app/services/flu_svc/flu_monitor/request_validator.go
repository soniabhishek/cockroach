package flu_monitor

import (
	"github.com/crowdflux/angel/app/plog"
	"net/http"
)

func shouldRetry(resp *http.Response) (*FluResponse, bool) {
	defer resp.Body.Close()
	fluResp := ParseFluResponse(resp)
	shouldCallBack := HttpCodeForCallback(fluResp.HttpStatusCode)
	plog.Trace("HTTPStatusCode: [", fluResp.HttpStatusCode, "] Should Call back: ", shouldCallBack)
	if shouldCallBack {
		return fluResp,true
	} else {
		//If any invalid flu response code is in our InvalidationCodeArray, then we log[ERROR] it
		for _, invalidFlu := range fluResp.Invalid_Flus {
			if IsValidInternalError(invalidFlu.Error) {
				return fluResp, false
			}
		}
	}
	return fluResp, true
}
