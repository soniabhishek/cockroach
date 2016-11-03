package work_flow_svc

import (
	"fmt"
	"testing"
	"time"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWorkFlowSvc_AddFLU(t *testing.T) {

	t.SkipNow()

	fluId := uuid.NewV4()

	workFlowSvc := &workFlowSvc{}

	completeHandler := func(flu models.FeedLineUnit) {
		fmt.Println("on complete handler called", flu)
		assert.Equal(t, fluId, flu.ID)
	}
	workFlowSvc.OnComplete(completeHandler)

	workFlowSvc.Start()

	if false {

		go func() {

			for {

				workFlowSvc.AddFLU(models.FeedLineUnit{
					ID:          fluId,
					ReferenceId: uuid.NewV4().String(),
					ProjectId:   uuid.FromStringOrNil("6b6e70de-7fa1-483d-a0eb-02a979e5bc3b"),
				})

				time.Sleep(time.Duration(50) * time.Millisecond)
			}
		}()

		go func() {

			for {

				workFlowSvc.AddFLU(models.FeedLineUnit{
					ID:          fluId,
					ReferenceId: uuid.NewV4().String(),
					ProjectId:   uuid.FromStringOrNil("6b6e70de-7fa1-483d-a0eb-02a979e5bc3b"),
				})

				time.Sleep(time.Duration(50) * time.Millisecond)
			}
		}()

		go func() {

			for {

				workFlowSvc.AddFLU(models.FeedLineUnit{
					ID:          fluId,
					ReferenceId: uuid.NewV4().String(),
					ProjectId:   uuid.FromStringOrNil("6b6e70de-7fa1-483d-a0eb-02a979e5bc3b"),
				})

				time.Sleep(time.Duration(500) * time.Millisecond)
			}
		}()
	}

	time.Sleep(time.Duration(100) * time.Second)

}
