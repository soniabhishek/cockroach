package crowdsourcing_step

import (
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
)

func New() *crowdSourcingStep {
	return &crowdSourcingStep{
		Step: step.New(),
	}
}
