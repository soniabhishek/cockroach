package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
)

var base = "https://load.playment.in"

//var base = "http://localhost:9000"
var getUsersUrl = "/api/post_processing_data?type=userBulk&limit="
var getMicroTaskUrl = base + "/api/micro_tasks/available"
var submitMissionUrl = base + "/api/mission_submissions"

type User struct {
	Id   string `json:"id"`
	Auth string `json:"auth"`
}

func main() {

	count := flag.Int("count", 10, "please provide the count")
	flag.Parse()
	Start(*count)
}

func Start(count int) {

	rand.Seed(100)

	users := GetUsers(count)

	fmt.Println("starting for", len(users))

	for _, usr := range users {

		go func() {
			for {
				fmt.Println("start a cycle")
				Sim(usr)
				time.Sleep(time.Duration(2) * time.Second)
			}
		}()
	}

	time.Sleep(time.Duration(500) * time.Minute)
}

func GetUsers(count int) []User {

	plog.Trace("getuser start")
	req, _ := http.NewRequest("GET", base+getUsersUrl+strconv.Itoa(count), nil)

	req.Header.Add("authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6IjQ1MmEzYTFkLTEwNzYtNDc3MC1hODg0LTZjNGU0MTExNzM1OSIsImlhdCI6MTQ2NzgzNjkxM30.jsqmb7z8EWlMcqUnUHXjliERTnrKgpEj5Hm3XHoHEPs")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if res == nil || res.Body == nil {
		log.Fatalln("res nil on get microtask list")
		plog.Trace("getmicrotasklist finish")
	}

	defer res.Body.Close()
	bty, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	type userListReq struct {
		Success bool   `json:"success"`
		Users   []User `json:"users"`
	}

	var userList userListReq

	err = json.Unmarshal(bty, &userList)
	if err != nil {
		log.Fatalln(err)
	}
	users := []User{}

	for _, usr := range userList.Users {
		usr.Auth = "Bearer " + usr.Auth
		users = append(users, usr)
	}
	return users
}

func Sim(user User) {

	microTasks := GetMicroTaskList(user)
	if !microTasks.Success || len(microTasks.MicroTask) == 0 {
		plog.Info("success false return")
		return
	}

	indexToPick := rand.Intn(len(microTasks.MicroTask))

	microTaskId := microTasks.MicroTask[indexToPick]["id"].(string)

	mission := StartMission(user, microTaskId)
	if mission == nil || len(mission.Questions) == 0 {
		plog.Trace("mission nil or questions empty")
		return
	}

	missionId := mission.Mission["id"].(string)
	questionSubmission := []models.JsonF{}

	for _, q := range mission.Questions {

		plog.Trace(q.String())

		qs := models.JsonF{
			"question_id": q["id"].(string),
			"answer": models.JsonF{
				"body": models.JsonF{
					"choice": "Flats",
				},
			},
		}
		questionSubmission = append(questionSubmission, qs)
	}

	t := 30 + rand.Intn(300)
	time.Sleep(time.Duration(t) * time.Second)

	SubmitMission(user, MissionSubmission{missionId, questionSubmission})
}

type MicroTaskList struct {
	MicroTask []models.JsonF `json:"micro_task"`
	Success   bool           `json:"success"`
}

func GetMicroTaskList(user User) *MicroTaskList {
	plog.Trace("getmicrotasklist start")
	req, _ := http.NewRequest("GET", getMicroTaskUrl, nil)

	req.Header.Add("authorization", user.Auth)
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if res == nil || res.Body == nil {
		log.Fatalln("res nil on get microtask list")
		plog.Trace("getmicrotasklist finish")
	}

	defer res.Body.Close()
	bty, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var microTaskList MicroTaskList

	err = json.Unmarshal(bty, &microTaskList)
	if err != nil {
		log.Fatalln(err)
	}

	plog.Trace("getmicrotasklist finish")
	return &microTaskList
}

func StartMission(user User, microTaskId string) *Mission {
	plog.Trace("StartMission start")

	bty := []byte(`{"mission":{"micro_task_id":"` + microTaskId + `"}}`)

	req, _ := http.NewRequest("POST", base+`/api/micro_tasks/`+microTaskId+`/missions`, bytes.NewBuffer(bty))

	req.Header.Add("authorization", user.Auth)
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if res == nil || res.Body == nil {
		log.Fatalln("empty response create mission")
	}

	if res.StatusCode != http.StatusOK {
		return nil
	}

	defer res.Body.Close()
	bty, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var mission Mission

	err = json.Unmarshal(bty, &mission)
	if err != nil {
		log.Fatalln(err)
	}
	plog.Trace("StartMission finish")

	return &mission
}

func SubmitMission(user User, msr MissionSubmission) {

	plog.Trace("SubmitMission start")

	reqBody := MissionSubmissionReq{
		MissionSubmission: msr,
	}
	plog.Trace("SubmitMission 1")

	reqBodyBty, _ := json.Marshal(reqBody)

	plog.Trace(string(reqBodyBty))

	req, _ := http.NewRequest("POST", submitMissionUrl, bytes.NewBuffer(reqBodyBty))

	req.Header.Add("authorization", user.Auth)
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	plog.Trace("SubmitMission 2")

	if err != nil {
		log.Fatalln(err)
	}

	plog.Trace("SubmitMission 3")

	if res == nil || res.Body == nil {
		log.Fatalln("empty response submit mission")
	}

	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	plog.Trace("SubmitMission finish")
}

type Mission struct {
	Mission   models.JsonF   `json:"mission"`
	Questions []models.JsonF `json:"questions"`
}

type MissionSubmissionReq struct {
	MissionSubmission MissionSubmission `json:"mission_submission"`
}

type MissionSubmission struct {
	MissionId          string         `json:"mission_id"`
	QuestionSubmission []models.JsonF `json:"question_submission"`
}
