package flu_svc

import (
	"gitlab.com/playment-main/support/app/models"
	"gitlab.com/playment-main/support/app/models/uuid"
)

type IFluService interface {
	AddFeedLineUnit(flu *models.FeedLineUnit) error
	SyncInputFeedLine() error
	GetFeedLineUnit(fluId uuid.UUID) (models.FeedLineUnit, error)
}

type IFluServiceExtended interface {
	IFluService
	IValidatorService
}

type IValidatorService interface {
	GetValidators(macroTaskId uuid.UUID, tag string) (fvs []models.FLUValidator, err error)
	SaveValidator(fv *models.FLUValidator) (err error)
}
