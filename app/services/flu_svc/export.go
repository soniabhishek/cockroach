package flu_svc

import (
	"fmt"
	"github.com/robfig/cron"
	"gitlab.com/playment-main/support/app/services/data_access_svc/repositories/feed_line_repo"
	"gitlab.com/playment-main/support/app/services/data_access_svc/repositories/flu_validator_repo"
	"gitlab.com/playment-main/support/app/services/data_access_svc/repositories/macro_task_repo"
)

func New() IFluService {
	return &fluService{
		fluRepo:          feed_line_repo.New(),
		fluValidatorRepo: flu_validator_repo.New(),
		macroTaskRepo:    macro_task_repo.New(),
	}
}

func NewWithExposedValidators() IFluServiceExtended {
	return &extendedFluService{
		fluService: fluService{
			fluRepo:          feed_line_repo.New(),
			fluValidatorRepo: flu_validator_repo.New(),
			macroTaskRepo:    macro_task_repo.New(),
		},
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
