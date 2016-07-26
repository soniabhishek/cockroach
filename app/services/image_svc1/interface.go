package image_svc1

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IImageService interface {
	BulkDownloadImages([]models.ImageDictionaryNew) (batchId uuid.UUID, err error)
	GetBatchesForMicroTask(microTaskId uuid.UUID) (batches []models.BatchProcess, err error)
	GetDownloadStatus(batchId uuid.UUID) (models.BatchProcess, error)
}
