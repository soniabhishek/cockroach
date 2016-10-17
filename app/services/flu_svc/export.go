package flu_svc

import (
	"fmt"

	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_validator"
	"github.com/crowdflux/angel/app/services/work_flow_svc"
	"strconv"
	"time"
)

func New() IFluService {
	return &fluService{
		fluRepo:      feed_line_repo.New(),
		fluValidator: flu_validator.New(),
		projectsRepo: projects_repo.New(),
		workFlowSvc:  work_flow_svc.StdWorkFlowSvc,
	}
}

type extendedFluService struct {
	fluService
	flu_validator.IFluValidatorService
}

func NewWithExposedValidators() IFluServiceExtended {
	return &extendedFluService{
		fluService: fluService{
			fluRepo:      feed_line_repo.New(),
			fluValidator: flu_validator.New(),
			projectsRepo: projects_repo.New(),
		},
		IFluValidatorService: flu_validator.New(),
	}
}

func StartFeedLineSync() {

	go func() {

		fSvc := New()

		intervalInSec, err := strconv.Atoi(config.INPUT_FEEDLINE_SYNC_TIME_PERIOD_SEC.Get())
		if err != nil {
			panic(err)
		}

		ticker := time.Tick(time.Duration(intervalInSec) * time.Second)

		plog.Info("Input Feedline", "started syncing at every "+strconv.Itoa(intervalInSec)+" seconds")

		for range ticker {
			err := fSvc.SyncInputFeedLine()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

}
