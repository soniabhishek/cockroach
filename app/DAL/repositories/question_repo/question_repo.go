package question_repo

import (
	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/models"
)

type questionRepo struct {
	db repositories.IDatabase
}

var _ IQuestionRepo = &questionRepo{}

func (r *questionRepo) Add(q models.Question) error {

	return r.db.Insert(&q)
}

func (r *questionRepo) Update(q models.Question) error {

	_, err := r.db.Update(&q)
	return err
}
