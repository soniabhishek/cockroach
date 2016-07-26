package flu_svc

import (
	"fmt"

	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_validator"
	"github.com/crowdflux/angel/app/services/work_flow_svc"
	"github.com/robfig/cron"
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

	fSvc := New()
	c := cron.New()

	syncFeedLine := func() {
		err := fSvc.SyncInputFeedLine()
		if err != nil {
			fmt.Println(err)
		}
	}

	c.AddFunc("*/10 * * * * *", syncFeedLine)
	c.Start()
}
