package crowdsourcing_step

import (
	"gitlab.com/playment-main/angel/app/DAL/clients"
	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
)

var StdCrowdSourcingStep = &crowdSourcingStep{
	Step:      step.New(),
	fluRepo:   feed_line_repo.New(),
	fluClient: clients.GetCrowdyClient(),
}

// Just a short form for above
var Std = StdCrowdSourcingStep
