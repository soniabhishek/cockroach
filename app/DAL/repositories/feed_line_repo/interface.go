package feed_line_repo

import (
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type IFluRepo interface {
	GetById(Id uuid.UUID) (models.FeedLineUnit, error)
	Save(feedLineUnit models.FeedLineUnit)
	BulkInsert(flus []models.FeedLineUnit) error
	Add(feedLineUnit models.FeedLineUnit) error
	Update(feedLineUnit models.FeedLineUnit) error
}

//--------------------------------------------------------------------------------//

type Mock struct{}

func (m *Mock) GetById(Id uuid.UUID) (models.FeedLineUnit, error) {
	return getMockFeedLineUnit(Id), nil
}

func (m *Mock) Save(ip models.FeedLineUnit) {
}

func (m *Mock) BulkInsert(flus []models.FeedLineUnit) error {
	return nil
}

func getMockFeedLineUnit(Id uuid.UUID) models.FeedLineUnit {

	return models.FeedLineUnit{
		ID:          Id,
		ReferenceId: "abcd",
		Data: models.JsonFake{
			"product_id": 123,
		},
		Tag: "dcba",
	}
}
