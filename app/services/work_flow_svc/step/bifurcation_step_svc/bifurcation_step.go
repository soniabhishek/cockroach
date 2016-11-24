package bifurcation_step_svc

import (
	"fmt"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
	"github.com/crowdflux/angel/app/services/work_flow_io_svc"
	"github.com/crowdflux/angel/app/services/work_flow_svc/step"
	"github.com/lib/pq"
	"time"
)

type bifurcationStep struct {
	step.Step
	fluRepo       feed_line_repo.IFluRepo
	stepConfigSvc work_flow_io_svc.IStepConfigSvc
}

const index string = "index"

func (b *bifurcationStep) processFlu(flu feed_line.FLU) {

	bifurcationConfig, err := b.stepConfigSvc.GetBifurcationStepConfig(flu.StepId)
	if err != nil {
		plog.Error("Bifurcation Step", err, "error getting step", "fluId: "+flu.ID.String())
		flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.Bifurcation, "Error getting bifurcation Config", flu.Redelivered())
		return
	}

	if flu.Build == nil {
		flu.Build = models.JsonF{}
	}

	if bifurcationConfig.Multiplication > 0 {

		for i := 0; i < bifurcationConfig.Multiplication; i++ {

			newFlu := flu.Copy()
			// Index starts from 1
			newFlu.Build[index] = i + 1

			// Created new Flus with masterId original Id
			if i > 0 {
				newFlu.ID = uuid.NewV4()
				newFlu.CreatedAt = pq.NullTime{time.Now(), true}
				newFlu.UpdatedAt = newFlu.CreatedAt
				newFlu.MasterId = flu.MasterId
				newFlu.IsMaster = false

				err := b.fluRepo.Add(newFlu.FeedLineUnit)
				if err != nil {
					flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.Bifurcation, "Error saving duplicate flu to db", flu.Redelivered())
					continue
				}
			} else {
				err := b.fluRepo.Update(newFlu.FeedLineUnit)
				if err != nil {
					flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.Bifurcation, "Error updating flu to db", flu.Redelivered())
					continue
				}
			}

			b.finishFlu(newFlu)
		}

	} else {
		plog.Error("Bifurcation Step", fmt.Errorf("", "Multiplication count not greater than 0"), "flu_id "+flu.ID.String())
		b.finishFlu(flu)
	}

	flu.ConfirmReceive()
}

func (c *bifurcationStep) finishFlu(flu feed_line.FLU) bool {
	c.OutQ.Push(flu)
	flu_logger_svc.LogStepExit(flu.FeedLineUnit, step_type.Bifurcation, flu.Redelivered())
	return true
}

func newStdBifurcation() *bifurcationStep {
	ts := &bifurcationStep{
		Step:          step.New(step_type.Bifurcation),
		fluRepo:       feed_line_repo.New(),
		stepConfigSvc: work_flow_io_svc.NewStepConfigService(),
	}

	ts.SetFluProcessor(ts.processFlu)
	return ts
}
