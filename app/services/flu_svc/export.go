package flu_svc

import (
	"fmt"

	"github.com/robfig/cron"
	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/DAL/repositories/macro_task_repo"
	"gitlab.com/playment-main/angel/app/services/flu_svc/flu_validator"
)

func New() IFluService {
	return &fluService{
		fluRepo:       feed_line_repo.New(),
		fluValidator:  flu_validator.New(),
		macroTaskRepo: macro_task_repo.New(),
	}
}

type extendedFluService struct {
	fluService
	flu_validator.IFluValidatorService
}

func NewWithExposedValidators() IFluServiceExtended {
	return &extendedFluService{
		fluService: fluService{
			fluRepo:       feed_line_repo.New(),
			fluValidator:  flu_validator.New(),
			macroTaskRepo: macro_task_repo.New(),
		},
		IFluValidatorService: flu_validator.New(),
	}
}

func StartFeedLineSync() {

	fSvc := New()
	c := cron.New()

	syncFeedLine := func() {
		err := fSvc.SyncInputFeedLine()
		if err != nil {
			fmt.Println(err)
		}
	}

	c.AddFunc("0/20 * * * * *", syncFeedLine)
	c.Start()
}
