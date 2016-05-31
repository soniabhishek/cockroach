package crowdsourcing_step

import (
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
)

var StdCrowdSourcingStep = &crowdSourcingStep{
	Step: step.New(),
}

// Just a short form for above
var Std = StdCrowdSourcingStep
