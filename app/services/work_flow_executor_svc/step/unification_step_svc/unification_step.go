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

	uuidsToSkip := []uuid.UUID{uuid.FromStringOrNil("179700d8-7ddd-48a9-b0e5-98c83b7eb144"), uuid.FromStringOrNil("6b278451-d350-4285-9ba9-849cb9a5bf2b"), uuid.FromStringOrNil("b9f71d53-fd8e-4acc-a980-d8f0a2b35667"), uuid.FromStringOrNil("4f494114-dbea-4c4f-a0ae-b29e7cd3ecad"), uuid.FromStringOrNil("15817696-c129-430f-a2bc-f7996fc861da"), uuid.FromStringOrNil("ac5d5b68-8b4b-4c76-9b21-ffa594d8dc16"), uuid.FromStringOrNil("5aeb0e8d-0e9a-43a6-a7e6-d4a86d4a7084"), uuid.FromStringOrNil("7db091b2-cac9-4d70-b086-137ab4b92255"), uuid.FromStringOrNil("593c1eb2-0a22-4ce8-b0b7-43a289efa6c8"), uuid.FromStringOrNil("ed5bac16-7824-4c90-b398-ac546e48a608"), uuid.FromStringOrNil("9db31b80-ca0e-469f-8904-9b82be606935"), uuid.FromStringOrNil("21a834af-289a-46af-80c1-9eff5617236a"), uuid.FromStringOrNil("b198817c-9118-46fe-9cb9-05ee632d4035"), uuid.FromStringOrNil("bace65ec-28a2-4606-bd20-658be604006d"), uuid.FromStringOrNil("0bd8dfb1-b6c8-4e42-a0d6-38ff0eff0750"), uuid.FromStringOrNil("1842c3a1-49e7-47d0-822c-d3ad9ce7637f")}
	for _, skipId := range uuidsToSkip {

		if skipId == flu.ID {
			fmt.Println("unification skipped", "fluid: "+flu.ID.String())
			flu.ConfirmReceive()
			return
		}
	}

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
