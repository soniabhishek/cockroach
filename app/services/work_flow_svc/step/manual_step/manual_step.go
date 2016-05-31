package manual_step

import (
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
)

type manualStep struct {
	step.Step
}

func processFlu(flu feed_line.FLU) {
	//flu.FeedLineUnit
}
