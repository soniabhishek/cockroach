package flu_svc

import (
	"fmt"

	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/DAL/repositories/flu_validator_repo"
	"gitlab.com/playment-main/angel/app/DAL/repositories/macro_task_repo"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type fluService struct {
	fluRepo          feed_line_repo.IFluRepo
	fluValidatorRepo flu_validator_repo.IFluValidatorRepo
	macroTaskRepo    macro_task_repo.IMacroTaskRepo
}

var _ IFluService = &fluService{}

func (i *fluService) AddFeedLineUnit(flu *models.FeedLineUnit) error {

	if flu.ReferenceId == "" {
		return ErrReferenceIdMissing
	}

	_, err := validateFlu(i.fluValidatorRepo, *flu)
	if err != nil {
		return err
	}

	err = checkMacroTaskExists(i.macroTaskRepo, flu.MacroTaskId)
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
