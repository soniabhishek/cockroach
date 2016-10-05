package flu_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_validator"
	"github.com/crowdflux/angel/app/services/work_flow_svc"
)

type fluService struct {
	fluRepo      feed_line_repo.IFluRepo
	fluValidator flu_validator.IFluValidatorService
	projectsRepo projects_repo.IProjectsRepo
	workFlowSvc  work_flow_svc.IWorkFlowSvc
}

var _ IFluService = &fluService{}

func (i *fluService) AddFeedLineUnit(flu *models.FeedLineUnit) error {

	_, err := i.fluValidator.Validate(*flu)
	if err != nil {
		return err
	}

	err = checkProjectExists(i.projectsRepo, flu.ProjectId)
	if err != nil {
		return err
	}

	fin := feed_line_repo.NewInputQueue()
	id, err := fin.Add(*flu)
	flu.ID = id
	if err != nil && err == feed_line_repo.ErrDuplicateReferenceId {
		err = ErrDuplicateReferenceId
	}

	return err
}

func (i *fluService) SyncInputFeedLine() error {

	fluInputQueue := feed_line_repo.NewInputQueue()

	flus, err := fluInputQueue.GetQueued()

	if err != nil {

		plog.Error("Error occured while getting data", err)
		return err
	}

	if len(flus) > 0 {

		err = i.fluRepo.BulkInsert(flus)

		if err != nil {
			plog.Error("Bulk insert failed", err)
			return err
		}

		// start adding to workFlowSvc in another go routine
		go func() {

			for _, flu := range flus {
				i.workFlowSvc.AddFLU(flu)
			}
		}()

		err = fluInputQueue.MarkFinished(flus)

		if err != nil {
			plog.Error("Changing queue status failed", err)
			return err
		}
		//plog.Info(len(flus), "flus processed")

	}

	return nil
}

func (i *fluService) GetFeedLineUnit(fluId uuid.UUID) (models.FeedLineUnit, error) {

	flu, err := i.fluRepo.GetById(fluId)
	if err != nil && err == feed_line_repo.ErrFLUNotFoundInInputQueue {
		err = ErrFluNotFound
	}
	return flu, err
}

//--------------------------------------------------------------------------------//
//CHECK PROJECT
//--------------------------------------------------------------------------------//

func checkProjectExists(r projects_repo.IProjectsRepo, mId uuid.UUID) error {
	_, err := r.GetById(mId)
	return err
}
