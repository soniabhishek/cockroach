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

	workFlowSvc.AddFLU(models.FeedLineUnit{
		ID: fluId,
	})

	time.Sleep(time.Duration(100) * time.Second)

}
