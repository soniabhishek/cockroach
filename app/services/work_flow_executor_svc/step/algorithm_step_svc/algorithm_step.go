package algorithm_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/clients"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
	"github.com/crowdflux/angel/app/services/work_flow_executor_svc/step"
	"github.com/crowdflux/angel/app/services/work_flow_svc"
	"math/rand"
	"time"
)

type algorithmStep struct {
	step.Step
	stepConfigSvc work_flow_svc.IStepConfigSvc
}

func (t *algorithmStep) processFlu(flu feed_line.FLU) {
	t.AddToBuffer(flu)

	plog.Info("algorithm Step flu reached", flu.ID)
	aStepConf, err := t.stepConfigSvc.GetAlgorithmStepConfig(flu.StepId)
	if err != nil {
		plog.Error("Algorithm step", err, plog.NewMessageWithParam("fluId: ", flu.ID.String()), plog.NewMessageWithParam("stepid: ", flu.StepId.String()), plog.NewMessageWithParam("Flu:", flu.FeedLineUnit))
		flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.Algorithm, "Algorithm Config Error", flu.Redelivered())
		return
	}

	textSting, ok := flu.Build[aStepConf.TextFieldKey].(string)
	if !ok {
		flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.Algorithm, "text field: "+aStepConf.TextFieldKey+" not present as string", flu.Redelivered())
		return
	}

	algoResult, err, success := clients.GetAbacusClient().Predict(textSting)
	if err != nil {
		plog.Error("Algorithm step", err, plog.NewMessageWithParam("fluId: ", flu.ID.String()), plog.NewMessageWithParam("flu", flu.FeedLineUnit))
	} else if success {
		flu.Build[aStepConf.AnswerKey] = algoResult
	}

	answerKeySuccess := aStepConf.AnswerKey + "_success"
	flu.Build[answerKeySuccess] = success

	if !success {

		t.finishFlu(flu)
		flu.ConfirmReceive()
		return
	}

	go func() {

		timeDiff := aStepConf.TimeDelayStop - aStepConf.TimeDelayStart

		if timeDiff <= 0 {
			plog.Info("Algostep", "timediff <= 0", "value: ", timeDiff, "fluId: "+flu.ID.String())
			timeDiff = 0
		}

		time.Sleep(time.Duration(int64((aStepConf.TimeDelayStart+timeDiff*rand.Float64())*60)) * time.Second)
		t.finishFlu(flu)
		flu.ConfirmReceive()
	}()

}

func (t *algorithmStep) finishFlu(flu feed_line.FLU) bool {

	err := t.RemoveFromBuffer(flu)
	if err != nil {
		plog.Trace("algorithm step", "flu not present in buffer")
	}
	t.OutQ.Push(flu)
	plog.Info("algorithm step out", flu.ID)
	flu_logger_svc.LogStepExit(flu.FeedLineUnit, step_type.Algorithm, flu.Redelivered())

	return true
}

func newStdPredictor() *algorithmStep {
	ts := &algorithmStep{
		Step:          step.New(step_type.Algorithm),
		stepConfigSvc: work_flow_svc.NewStepConfigService(),
	}

	ts.SetFluProcessor(ts.processFlu)
	return ts
}
