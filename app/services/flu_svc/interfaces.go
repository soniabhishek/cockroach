package flu_svc

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_validator"
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
