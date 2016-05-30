package question_repo

import (
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
	"gitlab.com/playment-main/angel/app/models"
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
