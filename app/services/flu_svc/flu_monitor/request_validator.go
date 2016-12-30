package flu_monitor

import (
	"github.com/crowdflux/angel/app/plog"
	"net/http"
)

func shouldRetry(resp *http.Response, retryLeft int) (*WebhookResponse, bool) {
	defer resp.Body.Close()
	fluResp := ParseFluResponse(resp)
	if retryLeft < 1 {
		return fluResp, false
	}
	shouldCallBack := HttpCodeForCallback(fluResp.HttpStatusCode)
	plog.Trace("HTTPStatusCode: [", fluResp.HttpStatusCode, "] Should Call back: ", shouldCallBack)
	if shouldCallBack {
		return fluResp, true
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
