package crowdsourcing_step_svc

import (
	"github.com/crowdflux/angel/app/DAL/clients"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
	"github.com/crowdflux/angel/app/services/work_flow_executor_svc/step"
	"github.com/crowdflux/angel/app/services/work_flow_svc"
)

type crowdSourcingStep struct {
	step.Step
	fluRepo       feed_line_repo.IFluRepo
	fluClient     crowdsourcingGatewayClient
	stepConfigSvc work_flow_svc.IStepConfigSvc
}

// Rename the interface later
type crowdsourcingGatewayClient interface {
	PushFLU(flu models.FeedLineUnit, microTaskId uuid.UUID) (bool, error)
}

func (c *crowdSourcingStep) processFlu(flu feed_line.FLU) {

	if step.IsSkipped(flu.ID) {
		flu.ConfirmReceive()
		return
	}

	c.AddToBuffer(flu)

	cc, err := c.stepConfigSvc.GetCrowdsourcingStepConfig(flu.StepId)
	if err != nil {
		plog.Error("crowdsourcing step", err, flu.ID.String())
		flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.CrowdSourcing, "ConfigNotFound", flu.Redelivered())
		return
	}

	_, err = c.fluClient.PushFLU(flu.FeedLineUnit, cc.MicroTaskId)
	if err != nil {
		plog.Error("crowdsourcing step", err, flu.ID.String())
		flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.CrowdSourcing, "crowdySendFailure", flu.Redelivered())
		return
	} else {
		flu.ConfirmReceive()
	}
}

func (c *crowdSourcingStep) finishFlu(flu feed_line.FLU) bool {

	err := c.RemoveFromBuffer(flu)
	if err != nil {
		plog.Trace("Crowdsourcing Step", err, "flu not present", flu.ID)
		//return false
	}
	c.OutQ.Push(flu)
	flu_logger_svc.LogStepExit(flu.FeedLineUnit, step_type.CrowdSourcing, flu.Redelivered())
	return true
}

func newCrowdSourcingStep() *crowdSourcingStep {

	cs := &crowdSourcingStep{
		Step:          step.New(step_type.CrowdSourcing),
		fluRepo:       feed_line_repo.New(),
		fluClient:     clients.GetCrowdyClient(),
		stepConfigSvc: work_flow_svc.NewStepConfigService(),
	}
	cs.Step.SetFluProcessor(cs.processFlu)
	return cs
}
