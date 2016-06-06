package work_flow_svc

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

func TestWorkFlowSvc_AddFLU(t *testing.T) {

	//t.SkipNow()

	fluId := uuid.NewV4()

	workFlowSvc := &workFlowSvc{}

	completeHandler := func(flu models.FeedLineUnit) {
		fmt.Println("on complete handler called", flu)
		assert.Equal(t, fluId, flu.ID)
	}
	workFlowSvc.OnComplete(completeHandler)

	workFlowSvc.Start()

	workFlowSvc.AddFLU(models.FeedLineUnit{
		ID: fluId,
	})

	time.Sleep(time.Duration(100) * time.Second)

}
