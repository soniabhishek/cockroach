package work_flow_svc

import (
	"fmt"

	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/services/flu_svc/flu_output"
)

type IWorkFlowSvc interface {
	AddFLU(models.FeedLineUnit)
}

func newStd() IWorkFlowSvc {

	fOut := flu_output.New()

	completeHandler := func(flu models.FeedLineUnit) {
		fmt.Println("on complete handler called", flu)
		fOut.AddToOutputQueue(flu)
	}

	workFlowSvc := &workFlowSvc{}

	workFlowSvc.OnComplete(completeHandler)

	workFlowSvc.Start()
	return workFlowSvc
}

var StdWorkFlowSvc = newStd()
