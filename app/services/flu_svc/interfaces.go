package flu_svc

import (
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/services/flu_svc/flu_validator"
)

type IFluService interface {
	AddFeedLineUnit(flu *models.FeedLineUnit) error
	SyncInputFeedLine() error
	GetFeedLineUnit(fluId uuid.UUID) (models.FeedLineUnit, error)
}

type IFluServiceExtended interface {
	IFluService
	flu_validator.IFluValidatorService
}
