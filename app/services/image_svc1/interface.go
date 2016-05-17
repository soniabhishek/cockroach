package image_svc1

import (
	"gitlab.com/playment-main/support/app/models"
	"gitlab.com/playment-main/support/app/models/uuid"
)

type IImageService interface {
	BulkDownloadImages([]models.ImageDictionary) (batchId uuid.UUID, err error)
	GetBatchesForMicroTask(microTaskId uuid.UUID) (batches []models.BatchProcess, err error)
	GetDownloadStatus(batchId uuid.UUID) (models.BatchProcess, error)
}
