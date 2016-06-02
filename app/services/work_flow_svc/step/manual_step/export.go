package manual_step

import (
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
	"gitlab.com/playment-main/angel/utilities"
)

var StdManualStep = &manualStep{
	Step: step.New(),
	id:   utilities.StartId,
}

// Just a short form for above
var Std = StdManualStep
