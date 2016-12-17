package flu_monitor

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/status_codes"
	"github.com/crowdflux/angel/app/plog"
	"net/http"
)

func addSendBackAuth(req *http.Request, fpsModel models.ProjectConfiguration, bodyJsonBytes []byte) {
	hmacKey := fpsModel.Options[HMAC_KEY]
	if hmacKey != nil {
		// ToDo add this when encrypted will be in DB
		//hmacKey, _ := utilities.Decrypt(hmacKey.(string))
		sig := hmac.New(sha256.New, []byte(hmacKey.(string)))
		sig.Write([]byte(string(bodyJsonBytes)))
		hmac := hex.EncodeToString(sig.Sum(nil))
		req.Header.Set(HMAC_HEADER_KEY, hmac)
		plog.Trace("HMAC", hmac)
	}
}

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
