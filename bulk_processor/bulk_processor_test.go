package bulk_processor_test

import (
	"fmt"
	"github.com/crowdflux/angel/bulk_processor"
	"testing"
	"time"
)

type TestClient struct {
	jobManager        *bulk_processor.JobManager
	internalFluPerSec int
	maxClientQps      int
	waitMiliSec       int
	name              string
}

func TestDispatcher_Start(t *testing.T) {

	clients := []TestClient{

		TestClient{
			jobManager:        bulk_processor.NewJobManager(1000, "1"),
			internalFluPerSec: 1000,
			maxClientQps:      1,
			waitMiliSec:       3,
			name:              "1",
		},
		TestClient{
			jobManager:        bulk_processor.NewJobManager(10, "2"),
			internalFluPerSec: 10,
			maxClientQps:      10,
			waitMiliSec:       3,
			name:              "2",
		},
	}

	dispatcher := bulk_processor.NewDispatcher(1)

	for _, c := range clients {
		dispatcher.AddJobManager(c.jobManager)
	}

	dispatcher.Start()

	for _, c := range clients {
		SendData(c)
	}

	time.Sleep(time.Minute * time.Duration(10))
}

func SendData(c TestClient) {

	go func() {
		ticker := time.Tick(time.Duration(1000/c.internalFluPerSec) * time.Millisecond)

		for {
			<-ticker
			c.jobManager.PushJob(bulk_processor.NewJob(func() {
				time.Sleep(time.Duration(c.waitMiliSec) * time.Millisecond)
				fmt.Println("finished " + c.name)
			}))
		}
	}()

}
