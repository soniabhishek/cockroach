package flu_output

import (
	"fmt"
	"sync"
	"time"
	"bytes"
	"net/http"
	"encoding/json"
	"gitlab.com/playment-main/angel/utilities"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/models/status_codes"
	"gitlab.com/playment-main/angel/app/DAL/repositories/flu_project_service"
)

var feedLinePipe = make(map[uuid.UUID]feedLineValue)
var retryCount = make(map[uuid.UUID]int)
var mutex = &sync.Mutex{}

type feedLineValue struct {
	insertionTime int64
	feedLine      []models.FeedLineUnit
}

type FluMonitor struct {

}

func (fm *FluMonitor) AddToOutputQueue(flu models.FeedLineUnit) error {

	feedLineArr := make([]models.FeedLineUnit, 1)
	feedLineArr[0] = flu
	return fm.AddManyToOutputQueue(feedLineArr)
}

func (fm *FluMonitor) AddManyToOutputQueue(fluBundle []models.FeedLineUnit) error {

	fmt.Println(feedLinePipe)
	mutex.Lock()
	for _, flu := range fluBundle {
		value, valuePresent := feedLinePipe[flu.ProjectID]
		if valuePresent {
			value = feedLineValue{utilities.TimeInMillis(), []models.FeedLineUnit{flu}}
		} else {
			value.feedLine = append(value.feedLine, flu)
		}
		feedLinePipe[flu.ProjectID] = value
	}
	mutex.Unlock()
	return nil
}

func checkupFeedLinePipe() {
	var projectIdsToSend = make([]uuid.UUID, 1)
	mutex.Lock()
	for projectId := range feedLinePipe {
		if IsEligibleForSendingBack(projectId) {
			projectIdsToSend = append(projectIdsToSend, projectId)
		}
	}
	mutex.Unlock()
	sendBackResp(projectIdsToSend)

}

func sendBackResp(projectIdsToSend []uuid.UUID) {

	retryIdsList := make([]uuid.UUID, 0)
	for _, projectId := range projectIdsToSend {
		flp, ok := feedLinePipe[projectId]
		if ok==false{
			continue
		}
		fluResp, status := sendBackToClient(projectId, flp.feedLine)
		if status == status_codes.Success {

			deleteFromFeedLinePipe(projectId)
			//TODO Should we write success logs?

		} else if status == status_codes.CallBackFailure && shouldRetryHttp(projectId){
			//not successful scenarios
			retryIdsList = append(retryIdsList, projectId)

		} else {
			//TODO Write TransactionErrorLog
			fmt.Println("Invalid FLU Response, ", fluResp)
			deleteFromFeedLinePipe(projectId)
		}
	}

	if len(retryIdsList) != 0 {
		time.Sleep(5000 * time.Millisecond) //TODO determine duration
		sendBackResp(retryIdsList)
	}
}

func sendBackToClient(projectId uuid.UUID, fluProjectResp []models.FeedLineUnit) (*models.Response, status_codes.StatusCode) {
	fpsRepo := flu_project_service.New()
	fpsModel , err := fpsRepo.Get(projectId)
	if utilities.IsValidError(err){
		fmt.Println(err)
		return &models.Response{}, status_codes.UnknownFailure
	}


	url := fpsModel.Url + "sdfadsf"
	//url := "http://localhost:8080/JServer/HelloServlet"
	fmt.Println("URL:>", url)

	jsonBytes, err := json.Marshal(fluProjectResp)
	if err != nil {
		//TODO check Error solid implementation
		fmt.Println(err)
		return &models.Response{}, status_codes.UnknownFailure
	}
	fmt.Println(string(jsonBytes))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//panic(err)
		return &models.Response{}, status_codes.UnknownFailure
	}

	fluResp, status := validationErrorCallback(resp)
	fluResp.FluStatusCode = status
	return fluResp, status
}

func validationErrorCallback(resp *http.Response) (*models.Response, status_codes.StatusCode) {
	defer resp.Body.Close()


	fluResp := utilities.ParseFluResponse(resp)
	fmt.Println(utilities.HttpCodeForCallback(fluResp.HttpStatusCode),fluResp.HttpStatusCode)
	if utilities.HttpCodeForCallback(fluResp.HttpStatusCode) {
		return fluResp, status_codes.CallBackFailure
	} else {
		//If any invalid flu response code is in our InvalidationCodeArray, then we log[ERROR] it
		for _, invalidFlu := range fluResp.Invalid_Flus {
			if utilities.IsValidInternalError(invalidFlu.Flu_Id) {
				return fluResp, status_codes.FluRespFailure
			}
		}
	}
	return fluResp, status_codes.Success
}

func IsEligibleForSendingBack(key uuid.UUID) bool{
	//TODO some threshold value and some configurable time
	flp, ok := feedLinePipe[key]
	if ok && (len(flp.feedLine) > 1000 || utilities.TimeDiff(false, flp.insertionTime) > 36000000) {
		return true
	}
	return false
}

/*-------------------------------------------------------------------------------------------------------------*/

func StartFluOutputTimer() {
	//Todo get scheduling value
	t := time.NewTicker(5 * time.Second)
	for _ = range t.C {
		checkupFeedLinePipe()
	}
}

func deleteFromFeedLinePipe(projectId uuid.UUID){
	mutex.Lock()
	delete(feedLinePipe, projectId)
	mutex.Unlock()
}

func shouldRetryHttp(projectId uuid.UUID) bool{
	prevRetryCnt, present := retryCount[projectId]
	if present == false || prevRetryCnt < 5{
		retryCount[projectId]++
		return true
	}else{
		delete(retryCount, projectId)
		return false
	}
}