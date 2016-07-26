package flu_output

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/project_configuration_repo"
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/status_codes"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/utilities"
)

var feedLinePipe = make(map[uuid.UUID]feedLineValue)
var retryCount = make(map[uuid.UUID]int)
var mutex = &sync.RWMutex{}
var dbLogger = feed_line_repo.StdLogger

var retryTimePeriod = time.Duration(utilities.GetInt(config.Get(config.RETRY_TIME_PERIOD))) * time.Millisecond

var defaultFluThresholdCount = utilities.GetInt(config.Get(config.DEFAULT_FLU_THRESHOLD_COUNT))
var fluThresholdDuration = int64(utilities.GetInt(config.Get(config.FLU_THRESHOLD_DURATION)))
var monitorTimePeriod = time.Duration(utilities.GetInt(config.Get(config.MONITOR_TIME_PERIOD))) * time.Millisecond
var retryThreshold = utilities.GetInt(config.Get(config.FLU_RETRY_THRESHOLD))

type feedLineValue struct {
	maxFluSize    int
	insertionTime int64
	feedLine      []models.FeedLineUnit
}

type FluMonitor struct {
}

type fluOutputStruct struct {
	ID          uuid.UUID   `json:"flu_id"`
	ReferenceId string      `json:"reference_id"`
	Tag         string      `json:"tag"`
	Status      string      `json:"status"`
	Result      interface{} `json:"result"`
}

func (fm *FluMonitor) AddToOutputQueue(flu models.FeedLineUnit) error {

	feedLineArr := make([]models.FeedLineUnit, 1)
	feedLineArr[0] = flu
	return fm.AddManyToOutputQueue(feedLineArr)
}

func (fm *FluMonitor) AddManyToOutputQueue(fluBundle []models.FeedLineUnit) error {

	plog.Info("FLu Monitor", fluBundle)

	mutex.Lock()
	for _, flu := range fluBundle {
		value, valuePresent := feedLinePipe[flu.ProjectId]
		if valuePresent == false {
			fpsRepo := project_configuration_repo.New()
			fpsModel, err := fpsRepo.Get(flu.ProjectId)
			if utilities.IsValidError(err) {
				plog.Error("DB Error:", err)
				return errors.New("No Project Configuration found for FluProject:" + flu.ProjectId.String())
			}
			maxFluCount := giveMaxFluCount(fpsModel)
			value = feedLineValue{maxFluCount, utilities.TimeInMillis(), []models.FeedLineUnit{flu}}
		} else {
			value.feedLine = append(value.feedLine, flu)
		}
		feedLinePipe[flu.ProjectId] = value
	}
	mutex.Unlock()
	return nil
}

func checkupFeedLinePipe() {

	plog.Trace("Flu output", "checkupFeedLinePipe")

	var projectIdsToSend = make([]uuid.UUID, 0)
	mutex.RLock()
	for projectId := range feedLinePipe {
		if IsEligibleForSendingBack(projectId) {
			projectIdsToSend = append(projectIdsToSend, projectId)
		}
	}
	mutex.RUnlock()
	if len(projectIdsToSend) > 0 {
		sendBackResp(projectIdsToSend)
	}

}

func sendBackResp(projectIdsToSend []uuid.UUID) {

	plog.Trace("Flu output", "sendBackResp", projectIdsToSend)

	retryIdsList := make([]uuid.UUID, 0)
	for _, projectId := range projectIdsToSend {
		flp, ok := feedLinePipe[projectId]
		if ok == false {
			continue
		}
		fluOutObj := getFluOutputObj(flp)

		fluResp, status := sendBackToClient(projectId, fluOutObj)
		if status == status_codes.Success {

			completedFLUs := deleteFromFeedLinePipe(projectId, fluOutObj)
			go putDbLog(completedFLUs, SUCCESS, *fluResp)

		} else if status == status_codes.CallBackFailure && shouldRetryHttp(projectId) {
			//not successful scenarios
			retryIdsList = append(retryIdsList, projectId)

		} else {
			completedFLUs := deleteFromFeedLinePipe(projectId, fluOutObj)
			go putDbLog(completedFLUs, "Invalid FLU Resp ", *fluResp)
		}
	}

	if len(retryIdsList) != 0 {
		time.Sleep(retryTimePeriod)
		sendBackResp(retryIdsList)
	}
}

func getFluOutputObj(flp feedLineValue) (fluOutputObj []fluOutputStruct) {
	flus := flp.feedLine
	limit := flp.maxFluSize
	if len(flp.feedLine) < flp.maxFluSize {
		limit = len(flp.feedLine)
	}
	plog.Info("SENDING FLUs COUNT: ", limit)
	for i := limit - 1; i >= 0; i-- {
		flu := flus[i]
		result, ok := flu.Build[RESULT]
		if !ok {
			result = models.JsonFake{}
		}

		fluOutputObj = append(fluOutputObj, fluOutputStruct{
			ID:          flu.ID,
			ReferenceId: flu.ReferenceId,
			Tag:         flu.Tag,
			Status:      STATUS_OK,
			Result:      result,
		})
	}
	return
}

func sendBackToClient(projectId uuid.UUID, fluProjectResp []fluOutputStruct) (*Response, status_codes.StatusCode) {

	if len(fluProjectResp) < 1 {
		return &Response{}, status_codes.NoFluToSend
	}

	plog.Info("Flu output", "sendBackToClient", projectId)

	fpsRepo := project_configuration_repo.New()
	fpsModel, err := fpsRepo.Get(projectId)
	if utilities.IsValidError(err) {
		plog.Error("DB Error:", err)
		return &Response{}, status_codes.UnknownFailure
	}

	url := fpsModel.PostBackUrl
	//url := "http://localhost:8080/JServer/HelloServlet"
	plog.Trace("URL:>", url, "|ID: ", projectId, "|Body:", fluProjectResp)

	sendResp := make(map[string][]fluOutputStruct)
	sendResp["feed_line_units"] = fluProjectResp
	jsonBytes, err := json.Marshal(sendResp)
	if err != nil {
		plog.Error("JSON Marshalling Error:", err)
		return &Response{}, status_codes.UnknownFailure
	}
	jsonBytes = utilities.ReplaceEscapeCharacters(jsonBytes)
	plog.Trace("Sending JSON:", string(jsonBytes))

	//fmt.Println(hex.EncodeToString(sig.Sum(nil)))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	req.Header.Set(CONTENT_TYPE, TYPE_JSON)

	for headerKey, headerVal := range fpsModel.Headers {
		req.Header.Set(headerKey, headerVal.(string))

	}
	addSendBackAuth(req, fpsModel, jsonBytes)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		plog.Error("HTTP Error:", err)
		return &Response{}, status_codes.UnknownFailure
	}

	fluResp, status := validationErrorCallback(resp)
	fluResp.FluStatusCode = status
	return fluResp, status
}

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

func validationErrorCallback(resp *http.Response) (*Response, status_codes.StatusCode) {
	defer resp.Body.Close()

	fluResp := ParseFluResponse(resp)
	shouldCallBack := HttpCodeForCallback(fluResp.HttpStatusCode)
	plog.Trace("HTTPStatusCode: [", fluResp.HttpStatusCode, "] Should Call back: ", shouldCallBack)
	if shouldCallBack {
		return fluResp, status_codes.CallBackFailure
	} else {
		//If any invalid flu response code is in our InvalidationCodeArray, then we log[ERROR] it
		for _, invalidFlu := range fluResp.Invalid_Flus {
			if IsValidInternalError(invalidFlu.Flu_Id) {
				return fluResp, status_codes.FluRespFailure
			}
		}
	}
	return fluResp, status_codes.Success
}

func IsEligibleForSendingBack(key uuid.UUID) bool {
	flp, ok := feedLinePipe[key]
	if !ok {
		return false
	}

	if len(flp.feedLine) < 1 {
		return false
	}

	if utilities.TimeDiff(false, flp.insertionTime) > fluThresholdDuration {
		return true
	}

	if len(flp.feedLine) > flp.maxFluSize {
		return true
	}

	return false
}

/*-------------------------------------------------------------------------------------------------------------*/

var startFluOnce sync.Once

func StartFluOutputTimer() {
	startFluOnce.Do(func() {
		plog.Trace("Flu output", monitorTimePeriod, "timer")

		t := time.NewTicker(monitorTimePeriod)
		go func() {
			for _ = range t.C {
				checkupFeedLinePipe()
			}
		}()
	})

}

func deleteFromFeedLinePipe(projectId uuid.UUID, fluOutputObj []fluOutputStruct) []models.FeedLineUnit {
	completedFLUs := make([]models.FeedLineUnit, 0)
	if len(fluOutputObj) < 1 {
		return completedFLUs
	}
	printFluBuff("BEFORE DELETION")
	mutex.Lock()
	flv, ok := feedLinePipe[projectId]
	if ok {
		for i := len(flv.feedLine) - 1; i >= 0; i-- {
			fl := flv.feedLine[i]
			// Condition to decide if current element has to be deleted:
			if didWeSendThis(fl, fluOutputObj) {

				completedFLUs = append(completedFLUs, flv.feedLine[i])

				flv.feedLine = append(flv.feedLine[:i], flv.feedLine[i+1:]...)

			}
		}
	}
	feedLinePipe[projectId] = flv
	mutex.Unlock()
	printFluBuff("AFTER DELETION")
	return completedFLUs
}

func didWeSendThis(fl models.FeedLineUnit, fluOutputObj []fluOutputStruct) bool {
	if fl.ID == uuid.Nil {
		return false
	}
	if len(fluOutputObj) > 0 {
		for i := len(fluOutputObj) - 1; i >= 0; i-- {
			if fluOutputObj[i].ID == uuid.Nil {
				continue
			}

			if fluOutputObj[i].ID == fl.ID {
				return true
			}
		}
	}
	return false
}

func printFluBuff(flag string) {
	if plog.IsTraceEnabled() {
		mutex.RLock()
		plog.Trace(flag, "OUTPUT FLU BUFF")
		for projectId := range feedLinePipe {
			plog.Trace("PROJ_ID:", projectId, "|BUFF-SIZE:", len(feedLinePipe[projectId].feedLine))
		}
		mutex.RUnlock()
	}
}

func shouldRetryHttp(projectId uuid.UUID) bool {
	prevRetryCnt, present := retryCount[projectId]
	if present == false || prevRetryCnt < retryThreshold {
		retryCount[projectId]++
		return true
	} else {
		delete(retryCount, projectId)
		return false
	}
}

func giveMaxFluCount(fpsModel models.ProjectConfiguration) int {
	val := fpsModel.Options[MAX_FLU_COUNT]
	if val == nil {
		return defaultFluThresholdCount
	}
	maxFluCount := utilities.GetInt(val.(string))
	if maxFluCount == 0 {
		maxFluCount = defaultFluThresholdCount
	}
	return maxFluCount
}
