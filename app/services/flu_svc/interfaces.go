package flu_svc

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_validator"
	"mime/multipart"
)

type IFluService interface {
	AddFeedLineUnit(flu *models.FeedLineUnit) error
	SyncInputFeedLine() error
	GetFeedLineUnit(uuid.UUID) (models.FeedLineUnit, error)
	CheckProjectExists(uuid.UUID) error
	GetUploadStatus(projectId string) (models.FluUploadStats, error)
	CsvCheckBasicValidation(file multipart.File, fileName string, projectId uuid.UUID) error
}

type IFluServiceExtended interface {
	IFluService
	flu_validator.IFluValidatorService
}
