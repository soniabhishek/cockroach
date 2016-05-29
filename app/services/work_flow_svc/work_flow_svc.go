package work_flow_svc

import (
	"fmt"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/counter"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/work_flow"
	"sync"
)

type workFlowSvc struct {
	work_flow.WorkFlow
	complete  OnCompleteHandler
	startOnce sync.Once
}

func (w *workFlowSvc) AddFLU(flu models.FeedLineUnit) {
	flu.Step = "workflowsvc"
	counter.Print(flu)
	w.InQ <- flu
}

func (w *workFlowSvc) Start() {

	//Executes only once, even if Start() is called multiple times
	w.startOnce.Do(func() {

		w.WorkFlow = work_flow.StdWorkFlow

		if w.complete != nil {
			startWorkflowSvc(w)
		} else {
			startWorkflowSvcNLog(w)
		}
	})
}

type OnCompleteHandler func(models.FeedLineUnit)

func (w *workFlowSvc) OnComplete(f OnCompleteHandler) {
	w.complete = f
}

func startWorkflowSvc(w *workFlowSvc) {
	go func() {
		for {
			select {
			case flu := <-w.OutQ:
				w.complete(flu)
			}
		}
	}()
}

func startWorkflowSvcNLog(w *workFlowSvc) {
	go func() {
		for {
			select {
			case flu := <-w.OutQ:
				fmt.Println(flu.ID)
			}
		}
	}()
}

func Start() {

	completeHandler := func(flu models.FeedLineUnit) {
		fmt.Println("on complete handler called", flu)
	}
	workFlowSvc := &workFlowSvc{}

	workFlowSvc.OnComplete(completeHandler)

	workFlowSvc.Start()
}
