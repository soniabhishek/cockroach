package flu_output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"gitlab.com/playment-main/angel/app/DAL/repositories/project_configuration_repo"
	"gitlab.com/playment-main/angel/app/config"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/status_codes"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/utilities"
)

var feedLinePipe = make(map[uuid.UUID]feedLineValue)
var retryCount = make(map[uuid.UUID]int)
var mutex = &sync.RWMutex{}

var retryTimePeriod = time.Duration(utilities.GetInt(config.Get(config.RETRY_TIME_PERIOD))) * time.Millisecond
var fluThresholdCount = utilities.GetInt(config.Get(config.FLU_THRESHOLD_COUNT))
var fluThresholdDuration = int64(utilities.GetInt(config.Get(config.FLU_THRESHOLD_DURATION)))
var monitorTimePeriod = time.Duration(utilities.GetInt(config.Get(config.MONITOR_TIME_PERIOD))) * time.Millisecond

type feedLineValue struct {
	insertionTime int64
	feedLine      []models.FeedLineUnit
}

type FluMonitor struct {
}

type fluOutputStruct struct {
	ID          uuid.UUID       `json:"flu_id"`
	ReferenceId string          `json:"reference_id"`
	Tag         string          `json:"tag"`
	Status      string          `json:"status"`
	Result      models.JsonFake `json:"results"`
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
		if valuePresent {
			value = feedLineValue{utilities.TimeInMillis(), []models.FeedLineUnit{flu}}
		} else {
			value.feedLine = append(value.feedLine, flu)
		}
		feedLinePipe[flu.ProjectId] = value
	}
	mutex.Unlock()
	return nil
}

func checkupFeedLinePipe() {

	plog.Info("Flu output", "checkupFeedLinePipe")

	var projectIdsToSend = make([]uuid.UUID, 1)
	mutex.RLock()
	for projectId := range feedLinePipe {
		if IsEligibleForSendingBack(projectId) {
			projectIdsToSend = append(projectIdsToSend, projectId)
		}
	}
	mutex.RUnlock()
	sendBackResp(projectIdsToSend)

}

func getFluOutputObj(flus []models.FeedLineUnit) (fluOutputObj []fluOutputStruct) {
	for _, flu := range flus {

		result := flu.Build["result"].(models.JsonFake)

		fluOutputObj = append(fluOutputObj, fluOutputStruct{
			ID:          flu.ID,
			ReferenceId: flu.ReferenceId,
			Tag:         flu.Tag,
			Status:      "COMPLETED",
			Result:      result,
		})
	}
	return
}

func sendBackResp(projectIdsToSend []uuid.UUID) {

	plog.Info("Flu output", "sendBackResp", projectIdsToSend)

	retryIdsList := make([]uuid.UUID, 0)
	for _, projectId := range projectIdsToSend {
		flp, ok := feedLinePipe[projectId]
		if ok == false {
			continue
		}
		fluOutObj := getFluOutputObj(flp.feedLine)

		fluResp, status := sendBackToClient(projectId, fluOutObj)
		if status == status_codes.Success {

			deleteFromFeedLinePipe(projectId)
			//TODO Should we write success logs?

		} else if status == status_codes.CallBackFailure && shouldRetryHttp(projectId) {
			//not successful scenarios
			retryIdsList = append(retryIdsList, projectId)

		} else {
			//TODO Write TransactionErrorLog
			fmt.Println("Invalid FLU Response, ", fluResp)
			deleteFromFeedLinePipe(projectId)
		}
	}

	if len(retryIdsList) != 0 {
		time.Sleep(retryTimePeriod * time.Millisecond)
		sendBackResp(retryIdsList)
	}
}

func sendBackToClient(projectId uuid.UUID, fluProjectResp []fluOutputStruct) (*Response, status_codes.StatusCode) {

	plog.Info("Flu output", "sendBackToClient", projectId)

	fpsRepo := project_configuration_repo.New()
	fpsModel, err := fpsRepo.Get(projectId)
	if utilities.IsValidError(err) {
		plog.Error("DB Error:", err)
		return &Response{}, status_codes.UnknownFailure
	}

	url := fpsModel.PostBackUrl
	//url := "http://localhost:8080/JServer/HelloServlet"
	fmt.Println("URL:>", url)

	jsonBytes, err := json.Marshal(fluProjectResp)
	if err != nil {
		plog.Error("JSON Marshalling Error:", err)
		return &Response{}, status_codes.UnknownFailure
	}
	fmt.Println(string(jsonBytes))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	for headerKey, headerVal := range fpsModel.Headers {
		req.Header.Set(headerKey, headerVal.(string))

	}

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

func validationErrorCallback(resp *http.Response) (*Response, status_codes.StatusCode) {
	defer resp.Body.Close()

	fluResp := ParseFluResponse(resp)
	fmt.Println(HttpCodeForCallback(fluResp.HttpStatusCode), fluResp.HttpStatusCode)
	if HttpCodeForCallback(fluResp.HttpStatusCode) {
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
	if ok && (len(flp.feedLine) > fluThresholdCount || utilities.TimeDiff(false, flp.insertionTime) > fluThresholdDuration) {
		return true
	}
	return false
}

/*-------------------------------------------------------------------------------------------------------------*/

var startFluOnce sync.Once

func StartFluOutputTimer() {
	startFluOnce.Do(func() {
		plog.Info("Flu output", monitorTimePeriod, "timer")

		t := time.NewTicker(monitorTimePeriod)
		go func() {
			for _ = range t.C {
				checkupFeedLinePipe()
			}
		}()
	})

}

func deleteFromFeedLinePipe(projectId uuid.UUID) {
	mutex.Lock()
	delete(feedLinePipe, projectId)
	mutex.Unlock()
}

func shouldRetryHttp(projectId uuid.UUID) bool {
	prevRetryCnt, present := retryCount[projectId]
	if present == false || prevRetryCnt < 5 {
		retryCount[projectId]++
		return true
	} else {
		delete(retryCount, projectId)
		return false
	}
}
