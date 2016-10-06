package work_flow_svc

import (
	"fmt"
	"sync"

	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/services/work_flow_svc/counter"
	"github.com/crowdflux/angel/app/services/work_flow_svc/work_flow"
)

type workFlowSvc struct {
	work_flow.WorkFlow
	complete  OnCompleteHandler
	startOnce sync.Once
}

func (w *workFlowSvc) AddFLU(flu models.FeedLineUnit) {
	counter.Print(feed_line.FLU{FeedLineUnit: flu}, "workflowsvc")
	w.InQ.Push(feed_line.FLU{FeedLineUnit: flu})
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
		for flu := range w.OutQ.Receiver() {
			w.complete(flu.FeedLineUnit)
			flu.ConfirmReceive()
		}
	}()
}

func startWorkflowSvcNLog(w *workFlowSvc) {
	go func() {
		for flu := range w.OutQ.Receiver() {
			fmt.Println(flu.ID)
			flu.ConfirmReceive()
		}
	}()
}
