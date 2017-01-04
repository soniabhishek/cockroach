package flu_monitor

import (
	"github.com/crowdflux/angel/app/plog"
	"net/http"
)

func shouldRetry(fluResp *WebhookResponse, retryLeft int) bool {
	if retryLeft < 1 || fluResp.HttpStatusCode == http.StatusOK {
		return false
	}
	shouldCallBack := HttpCodeForCallback(fluResp.HttpStatusCode)
	plog.Trace("HTTPStatusCode: [", fluResp.HttpStatusCode, "] Should Call back: ", shouldCallBack)
	if shouldCallBack {
		return true
	} else {
		//If any invalid flu response code is in our InvalidationCodeArray, then we log[ERROR] it
		for _, invalidFlu := range fluResp.Invalid_Flus {
			if IsValidInternalError(invalidFlu.Error) {
				return false
			}
		}
	}
	return false
}
