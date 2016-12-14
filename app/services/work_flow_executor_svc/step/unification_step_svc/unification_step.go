package unification_step_svc

import (
	"fmt"
	"github.com/crowdflux/angel/app/DAL/feed_line"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_logger_svc"
	"github.com/crowdflux/angel/app/services/work_flow_executor_svc/step"
	"github.com/crowdflux/angel/app/services/work_flow_svc"
)

type unificationStep struct {
	step.Step
	stepConfigSvc work_flow_svc.IStepConfigSvc
	fluRepo       feed_line_repo.IFluRepo
	fluCounter    fluCounter
}

func (u *unificationStep) processFlu(flu feed_line.FLU) {

	unificationConfig, err := u.stepConfigSvc.GetUnificationStepConfig(flu.StepId)
	if err != nil {
		plog.Error("Unification Step", err, "error getting step", "fluId: "+flu.ID.String())
		flu_logger_svc.LogStepError(flu.FeedLineUnit, step_type.Unification, "Error getting unification Config", flu.Redelivered())
		return
	}

	u.fluCounter.UpdateCount(flu)

	// If count is less then awaited count
	// then add it to the counter & don't finish it

	if u.fluCounter.GetCount(flu) >= unificationConfig.Multiplication {

		waitingFlus := u.fluCounter.Get(flu)

		var masterFlu feed_line.FLU

		for _, wFlu := range waitingFlus {
			flu.Build.Merge(wFlu.Build)

			if wFlu.IsMaster {
				masterFlu = wFlu
			} else {
				wFlu.IsActive = false
				err := u.fluRepo.Update(wFlu.FeedLineUnit)
				if err != nil {
					plog.Error("Unification Step", err, "error updating flu", "fluid "+flu.ID.String(), "masterid "+flu.MasterId.String())
					return
				}
			}
		}

		// master flu not found
		if masterFlu.ID == uuid.Nil {

			u.finishFlu(flu)
		} else {

			// if master flu found
			// then send it forward (reason not clear)
			masterFlu.Build = flu.Build.Copy()
			flu.IsActive = false

			err := u.fluRepo.Update(flu.FeedLineUnit)
			if err != nil {
				plog.Error("Unification Step", err, "error updating flu", "fluid "+flu.ID.String(), "masterid "+flu.MasterId.String())
				return
			}

			err = u.fluRepo.Update(masterFlu.FeedLineUnit)
			if err != nil {
				plog.Error("Unification Step", err, "error updating masterflu", "fluid "+masterFlu.ID.String(), "masterid "+flu.MasterId.String())
				return
			}

			u.finishFlu(masterFlu)
		}

		// clear counter
		u.fluCounter.Clear(flu)
		flu.ConfirmReceive()
	}
}

func (u *unificationStep) finishFlu(flu feed_line.FLU) bool {
	u.OutQ.Push(flu)
	flu_logger_svc.LogStepExit(flu.FeedLineUnit, step_type.Unification, flu.Redelivered())
	return true
}

func newStdUnificationStep() *unificationStep {
	ts := &unificationStep{
		Step:          step.New(step_type.Unification),
		fluRepo:       feed_line_repo.New(),
		stepConfigSvc: work_flow_svc.NewStepConfigService(),
		fluCounter:    newFluCounter(),
	}

	ts.SetFluProcessor(ts.processFlu)
	return ts
}
