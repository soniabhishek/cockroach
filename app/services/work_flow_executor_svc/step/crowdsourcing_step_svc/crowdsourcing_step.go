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

	if flu.ID == uuid.FromStringOrNil("7e764277-535e-4244-9f66-9d703b7540f2") ||
		flu.ID == uuid.FromStringOrNil("7c79cc52-a23d-4f09-ac3d-5d04d5bc74a2") ||
		flu.ID == uuid.FromStringOrNil("0837e141-65da-485c-a190-3bce180542c0") ||
		flu.ID == uuid.FromStringOrNil("8d723929-dc25-4108-8fc8-20ca0ca8ce3f") ||
		flu.ID == uuid.FromStringOrNil("cb03833e-62b1-4f10-980a-b951109a2c5b") ||
		flu.ID == uuid.FromStringOrNil("415d0fa3-ed53-4cbc-86eb-e2306c52cca6") ||
		flu.ID == uuid.FromStringOrNil("e575a788-315f-45a0-85cd-31a747921f55") ||
		flu.ID == uuid.FromStringOrNil("a31d360e-aa21-4b8c-93c0-34553fe6c5df") ||
		flu.ID == uuid.FromStringOrNil("ac30c3e5-2b85-4f52-851c-bf083875af9f") ||
		flu.ID == uuid.FromStringOrNil("0409cc0a-e046-40a8-a1ab-ddee4fad842d") ||
		flu.ID == uuid.FromStringOrNil("f98692f5-912e-4aed-90c7-dff255be7888") ||
		flu.ID == uuid.FromStringOrNil("b7d5d34d-be6f-4dd7-b5f4-1be00ae8116c") ||
		flu.ID == uuid.FromStringOrNil("864d6201-1d4d-4414-a5b4-cb264cd2f915") ||
		flu.ID == uuid.FromStringOrNil("cf364845-4a92-4193-9e5c-3496607c4da7") {
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
