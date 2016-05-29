package work_flow_svc

import "gitlab.com/playment-main/angel/app/models"

type IWorkFlowSvc interface {
	AddFLU(models.FeedLineUnit)
}
