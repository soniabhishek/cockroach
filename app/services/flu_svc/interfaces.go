package flu_svc

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/flu_upload_status"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_validator"
	"mime/multipart"
	"os"
)

type IFluService interface {
	AddFeedLineUnit(flu *models.FeedLineUnit) error
	SyncInputFeedLine() error
	GetFeedLineUnit(uuid.UUID) (models.FeedLineUnit, error)
	BulkAddFeedLineUnit(multipart.File, *os.File, string, uuid.UUID)
	CheckProjectExists(uuid.UUID) error
	CheckCsvUploaded(projectId string) (bool, error)
	UpdateUploadStatus(projectId uuid.UUID, status flu_upload_status.FluUploadStatus)
	GetUploadStatus(projectId string) (models.FluUploadStats, error)
}

type IFluServiceExtended interface {
	IFluService
	flu_validator.IFluValidatorService
}
