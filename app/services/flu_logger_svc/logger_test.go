package flu_logger_svc

import (
	"fmt"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"testing"
	"time"
)

type fluLoggerMock struct {
}

var _ feed_line_repo.IFluLogger = &fluLoggerMock{}

func (fl *fluLoggerMock) Log(logs []models.FeedLineLog) error {
	fmt.Println(len(logs))
	return nil
}

func TestLogStepEntry(t *testing.T) {

	t.Skipf("integration test")

	mockDbLogSyncer(&fluLoggerMock{})

	if true {

		go func() {

			for {
				LogStepEntry(models.FeedLineUnit{ID: uuid.NewV4()}, step_type.Manual, false)
			}

		}()
	}

	time.Sleep(time.Duration(100) * time.Second)
}
