package work_flow_svc

import (
	"fmt"
	"testing"
	"time"

	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

func completeHandler(flu models.FeedLineUnit) {
	fmt.Println("on complete handler called", flu)
}

func TestWorkFlowSvc_AddFLU(t *testing.T) {

	//t.SkipNow()

	fluId := uuid.NewV4()

	workFlowSvc := &workFlowSvc{}

	workFlowSvc.OnComplete(completeHandler)

	workFlowSvc.Start()

	workFlowSvc.AddFLU(models.FeedLineUnit{
		ID: fluId,
	})

	workFlowSvc.AddFLU(models.FeedLineUnit{
		ID: uuid.NewV4(),
	})
	workFlowSvc.AddFLU(models.FeedLineUnit{
		ID: uuid.NewV4(),
	})
	workFlowSvc.AddFLU(models.FeedLineUnit{
		ID: uuid.NewV4(),
	})
	time.Sleep(time.Duration(100) * time.Second)

}
