package crowdsourcing_step

import (
	"github.com/crowdflux/angel/app/DAL/clients"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
)

var StdCrowdSourcingStep = &crowdSourcingStep{
	Step:      step.New(),
	fluRepo:   feed_line_repo.New(),
	fluClient: clients.GetCrowdyClient(),
}

// Just a short form for above
var Std = StdCrowdSourcingStep
