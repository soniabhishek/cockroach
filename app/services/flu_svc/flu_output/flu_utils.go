package flu_output

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"database/sql"

	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/status_codes"
	"gitlab.com/playment-main/angel/app/models/uuid"
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

		/*TODO How to handle this error.
		errors.New("asfdsgdf")

		switch err.(type) {
		case :

		}*/

		plog.Error("Response Parsing Error: ", err)
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

func putDbLog(completedFLUs []models.FeedLineUnit, message string, resp Response) {
	dbLogArr := make([]models.FeedLineLog, len(completedFLUs))
	jsObj := models.JsonFake{}
	jsonBytes, _ := json.Marshal(resp)
	jsObj.Scan(string(jsonBytes))
	for i, fl := range completedFLUs {
		dbLog := models.FeedLineLog{
			//ID         int            `db:"id" json:"id" bson:"_id"`
			FluId:      fl.ID,
			Message:    sql.NullString{message, true},
			MetaData:   jsObj,
			StepType:   sql.NullInt64{int64(12), true},
			StepEntry:  sql.NullBool{true, true},
			StepExit:   sql.NullBool{true, true},
			StepId:     fl.StepId,
			WorkFlowId: uuid.UUID{},
			CreatedAt:  fl.CreatedAt,
		}
		dbLogArr[i] = dbLog
	}
	dbLogger.Log(dbLogArr)
}
