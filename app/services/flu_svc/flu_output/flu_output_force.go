package flu_output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/project_configuration_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/status_codes"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/utilities"
	"net/http"
	"time"
)

func ForceSendBackInParallel(stepId uuid.UUID, projectId uuid.UUID, concurrency int) {

	config, err := project_configuration_repo.New().Get(projectId)
	if err != nil {
		panic(err)
	}

	flus, err := feed_line_repo.New().GetFlusNotSent(stepId)
	if err != nil {
		panic(err)
	}

	fmt.Println("starting for flus: ", len(flus))

	var channelSize = concurrency

	c := make(chan int, channelSize)

	start := time.Now()

	for i := 0; i < len(flus)+channelSize; i++ {

		if i < len(flus) {
			go sendFluBack(config, flus[i], c)
		}
		if i >= channelSize {
			<-c
			fmt.Println(i-channelSize, time.Since(start), "fluID: "+flus[i-channelSize].ID.String())
		}
	}
}

func ForceSendBackInQps(stepId uuid.UUID, projectId uuid.UUID, qps int) {

	config, err := project_configuration_repo.New().Get(projectId)
	if err != nil {
		panic(err)
	}

	flus, err := feed_line_repo.New().GetFlusNotSent(stepId)
	if err != nil {
		panic(err)
	}

	fmt.Println("starting for flus: ", len(flus))

	start := time.Now()

	ch := make(chan int, 10000)

	go func() {
		k := 0
		for range ch {
			fmt.Println(time.Since(start), k)
			k++
		}
	}()

	for i := 0; i < len(flus); {

		for j := 0; j < qps; j++ {
			go sendFluBack(config, flus[i], ch)
		}
		i += qps
		time.Sleep(time.Duration(1) * time.Second)

	}
	time.Sleep(time.Duration(10) * time.Minute)
}

func sendFluBack(config models.ProjectConfiguration, flu models.FeedLineUnit, chn chan int) {
	defer func() {
		chn <- 1
	}()
	fluOutputStrct := fluOutputStruct{
		ID:          flu.ID,
		ReferenceId: flu.ReferenceId,
		Tag:         flu.Tag,
		Status:      STATUS_OK,
		Result:      flu.Build[RESULT],
	}

	fluResp, status := sendBackToClientCustom(config, []fluOutputStruct{fluOutputStrct})
	if status == status_codes.Success {

		putDbLog([]models.FeedLineUnit{flu}, SUCCESS, *fluResp)

	} else if status == status_codes.CallBackFailure {

		fmt.Println("failure", fluResp)

		//not successful scenarios
		//putDbLog([]models.FeedLineUnit{flu}, "ERROR", *fluResp)
	} else {
		fmt.Println("failure", fluResp)

		//putDbLog([]models.FeedLineUnit{flu}, "ERROR", *fluResp)
	}
}

func sendBackToClientCustom(fpsModel models.ProjectConfiguration, fluProjectResp []fluOutputStruct) (*Response, status_codes.StatusCode) {

	if len(fluProjectResp) < 1 {
		return &Response{}, status_codes.NoFluToSend
	}

	url := fpsModel.PostBackUrl
	//url := "http://localhost:8080/JServer/HelloServlet"
	plog.Trace("URL:>", url, "|ID: ", fpsModel.ProjectId, "|Body:", fluProjectResp)

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
