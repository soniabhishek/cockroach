package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"
)

var client *http.Client

func mainLoadTest() {

	client = &http.Client{
		Transport: &http.Transport{
			// just what default client uses
			Proxy: http.ProxyFromEnvironment,
			// this leads to more stable numbers
			MaxIdleConnsPerHost: 4 * runtime.GOMAXPROCS(0),
		},
	}

	var size = 400

	var buff chan int
	buff = make(chan int, size)

	for i := 0; i < size; i++ {
		buff <- 1
	}

	ticker := time.Tick(time.Duration(1) * time.Second)

	p := parallelSender{buff}
	//
	i := 0
	j := 0

	if true {

		// Send concurrently one by one
		go func() {
			for {
				select {
				case <-buff:
					i++
					go p.sendGetMicroTasks()
				case <-ticker:
					fmt.Println(i - j)
					j = i
				}
			}
		}()
	} else {

		// Send concurrently at every ticker tick
		go func() {
			for {
				select {
				case <-ticker:
					p := parallelSender{buff}
					for k := 0; k < size; k++ {
						go p.sendGetMicroTasks()
					}
				}
			}
		}()

	}

	time.Sleep(time.Duration(10) * time.Minute)

}

type availableMicroTaskReq struct {
}

var authkey = "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6IjQ1MmEzYTFkLTEwNzYtNDc3MC1hODg0LTZjNGU0MTExNzM1OSIsImlhdCI6MTQ2NTUxMjU0Nn0.mCzvkuHHMR6bAfCimL_WVaNhmH8sadWHcjdzkb1WdfY"

type parallelSender struct {
	Buff chan int
}

func (p *parallelSender) sendGetMicroTasks() {

	url := "http://localhost:9000/api/micro_tasks/available"

	req, _ := http.NewRequest("GET", url, nil)

	s := p.sendReq(req)
	fmt.Println(s)

}

func (p *parallelSender) sendCreateMission() {

	url := "http://localhost:9000/api/micro_tasks/62b6f8f6-5d11-450e-87d6-78aaaac28113/missions"

	bty := []byte(`{"mission":{"micro_task_id":"381b6c6f-4c32-439c-8770-11cabe3358a4"}}`)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bty))

	p.sendReq(req)

}

func (p *parallelSender) sendGetUsers() {
	url := "http://localhost:9000/api/users/452a3a1d-1076-4770-a884-6c4e41117359"

	req, _ := http.NewRequest("GET", url, nil)

	p.sendReq(req)
}

func (p *parallelSender) redeemCoupon() {
	url := "http://localhost:9000/api/coupon_redemptions"

	bty := []byte(`{"phone":"9632243142","email":"himanshu@playment.in","coupon_redemption":{"946ebf65-2832-4344-ae17-da2a8a2d76f0":1}}`)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bty))

	s := p.sendReq(req)

	if s != `{"success":false,"error":"Too many requests ascd"}` {
		fmt.Println(s)
	}
}

func (p *parallelSender) sendReq(req *http.Request) string {
	req.Header.Add("authorization", authkey)
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	if res == nil || res.Body == nil {
		p.Buff <- 1
		return ""
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	p.Buff <- 1
	return string(body)
}

func (p *parallelSender) sendNForget(req *http.Request) {
	req.Header.Add("authorization", authkey)
	req.Header.Add("content-type", "application/json")

	_, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	p.Buff <- 1
	return
}
