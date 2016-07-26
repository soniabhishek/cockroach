package question_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models"
)

type IQuestionRepo interface {
	Add(models.Question) error
	Update(models.Question) error
}

func New() IQuestionRepo {
	return &questionRepo{
		db: postgres.GetPostgresClient(),
	}
}
