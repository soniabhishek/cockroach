package bulk_processor_test

import (
	"testing"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_monitor/bulk_processor"
	"time"
	"fmt"
)

type TestClient struct {
	wm                bulk_processor.WorkerManager
	internalFluPerSec int
	maxClientQps      int
	waitMiliSec       int
	name              string
}

func TestDispatcher_Start(t *testing.T) {

	clients := []TestClient{
		{
			wm:                bulk_processor.NewWorkerManager(2),
			internalFluPerSec: 1,
			maxClientQps:      2,
			waitMiliSec:       1,
			name:              "C 1-10",
		},
		{
			wm:                bulk_processor.NewWorkerManager(1),
			internalFluPerSec: 1,
			maxClientQps:      1,
			waitMiliSec:       1,
			name:              "C 10-1",

		},
		//{
		//	wm:                bulk_processor.NewWorkerManager(5),
		//	internalFluPerSec: 5,
		//	maxClientQps:      5,
		//	name : "C 5-5",
		//
		//},
	}

	dispatcher := bulk_processor.NewDispatcher(2)

	for _, c := range clients {
		dispatcher.AddWorkerManager(c.wm)
	}

	dispatcher.Start()

	for _, c := range clients {
		go SendData(c)
	}

	time.Sleep(time.Minute * time.Duration(10))
}

func SendData(c TestClient) {

	ticker := time.Tick(time.Duration(1000 / c.internalFluPerSec) * time.Millisecond)

	for {
		<-ticker
		c.wm.PushJob(bulk_processor.NewJob(func() {
			fmt.Println(c.name)
		}))
	}
}
