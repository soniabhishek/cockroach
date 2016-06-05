package crowdsourcing_step

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/DAL/repositories/question_repo"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
	"testing"
)

type fakeQuestionRepo struct {
	questions map[uuid.UUID]models.Question
}

var _ question_repo.IQuestionRepo = &fakeQuestionRepo{}

func (f *fakeQuestionRepo) Add(q models.Question) error {
	f.questions[q.ID] = q
	return nil
}
func (f *fakeQuestionRepo) Update(q models.Question) error {
	f.questions[q.ID] = q
	return nil
}

//--------------------------------------------------------------------------------//

type fakeFluRepo struct {
	flus map[uuid.UUID]models.FeedLineUnit
}

var _ feed_line_repo.IFluRepo = &fakeFluRepo{}

func (m *fakeFluRepo) GetById(id uuid.UUID) (models.FeedLineUnit, error) {
	if flu, ok := m.flus[id]; ok {
		return flu, nil
	}
	return models.FeedLineUnit{}, errors.New("flu not found")
}

func (m *fakeFluRepo) Save(flu models.FeedLineUnit) {
	m.flus[flu.ID] = flu
}

func (m *fakeFluRepo) BulkInsert(flus []models.FeedLineUnit) error {
	for _, flu := range flus {
		m.Save(flu)
	}
	return nil
}
func (m *fakeFluRepo) Add(flu models.FeedLineUnit) error {
	m.flus[flu.ID] = flu
	return nil
}
func (m *fakeFluRepo) Update(flu models.FeedLineUnit) error {
	m.flus[flu.ID] = flu
	return nil
}

//--------------------------------------------------------------------------------//

func Test(t *testing.T) {

	fluRepo := fakeFluRepo{make(map[uuid.UUID]models.FeedLineUnit)}
	_ = fakeQuestionRepo{make(map[uuid.UUID]models.Question)}

	fluId := uuid.NewV4()

	flu := feed_line.FLU{
		FeedLineUnit: models.FeedLineUnit{
			ID:          fluId,
			ReferenceId: "PayFlip123",
			Tag:         "Brand",
			Data: models.JsonFake{
				"brand":  "Sony",
				"image1": "http://sxomeimaghe.com/some.jpeg",
			},
		},
	}

	fluRepo.Save(flu.FeedLineUnit)

	cs := crowdSourcingStep{
		Step:    step.New(),
		fluRepo: feed_line_repo.New(),
	}

	cs.InQ <- flu

	fluNew := <-cs.OutQ

	assert.EqualValues(t, flu.ID, fluNew.ID)
}
