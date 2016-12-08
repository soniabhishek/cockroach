package work_flow_executor_svc

import (
	"fmt"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_output"
)

type IWorkFlowSvc interface {
	AddFLU(models.FeedLineUnit)
}

func newStd() IWorkFlowSvc {

	fOut := flu_output.New()

	completeHandler := func(flu models.FeedLineUnit) {
		fmt.Println("on complete handler called", flu.ID)
		fOut.AddToOutputQueue(flu)
	}

	workFlowSvc := &workFlowExecutorSvc{}

	workFlowSvc.OnComplete(completeHandler)

	workFlowSvc.Start()
	return workFlowSvc
}

var StdWorkFlowSvc = newStd()
