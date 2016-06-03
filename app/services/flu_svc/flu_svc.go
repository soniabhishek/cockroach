package flu_svc

import (
	"fmt"

	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/DAL/repositories/projects_repo"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/services/flu_svc/flu_validator"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc"
)

type fluService struct {
	fluRepo      feed_line_repo.IFluRepo
	fluValidator flu_validator.IFluValidatorService
	projectsRepo projects_repo.IProjectsRepo
	workFlowSvc  work_flow_svc.IWorkFlowSvc
}

var _ IFluService = &fluService{}

func (i *fluService) AddFeedLineUnit(flu *models.FeedLineUnit) error {

	if flu.ReferenceId == "" {
		return ErrReferenceIdMissing
	}

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

		fmt.Println("Error occured while getting data", err)
		return err
	}

	if len(flus) > 0 {

		err = i.fluRepo.BulkInsert(flus)

		if err != nil {
			fmt.Println("Bulk insert failed", err)
			return err
		}

		err = fluInputQueue.MarkFinished()

		if err != nil {
			fmt.Println("Changing queue status failed")
			return err
		}
		fmt.Println(len(flus), "flus processed")

	}

	return nil
}

func (i *fluService) GetFeedLineUnit(fluId uuid.UUID) (models.FeedLineUnit, error) {
	fin := feed_line_repo.NewInputQueue()
	flu, err := fin.Get(fluId)
	if err != nil && err == feed_line_repo.ErrFLUNotFoundInInputQueue {
		err = ErrFluNotFound
	}
	return flu, err
}

//--------------------------------------------------------------------------------//
//CHECK PROJECT
//--------------------------------------------------------------------------------//

func checkProjectExists(r projects_repo.IProjectsRepo, mId uuid.UUID) error {
	_, err := r.Get(mId)
	return err
}
