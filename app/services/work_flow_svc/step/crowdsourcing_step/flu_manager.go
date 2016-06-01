package crowdsourcing_step

import (
	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/DAL/repositories/question_repo"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
)

type fluManager struct {
	QuestionRepo question_repo.IQuestionRepo
	FeedlineRepo feed_line_repo.IFluRepo
}

func (f *fluManager) DeriveQuestionBodyFromFlu(flu feed_line.FLU) models.JsonFake {
	return flu.Data
}

func newFluManager() fluManager {
	return fluManager{
		QuestionRepo: question_repo.New(),
		FeedlineRepo: feed_line_repo.New(),
	}
}
