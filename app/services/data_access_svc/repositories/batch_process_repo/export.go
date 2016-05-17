package batch_process_repo

import (
	"gitlab.com/playment-main/support/app/models"
	"gitlab.com/playment-main/support/app/models/uuid"
)

type IBatchProcessRepo interface {
	Get(uuid.UUID) (models.BatchProcess, error)
	Save(models.BatchProcess) error
}

//func New() IBatchProcessRepo {
//	return IBatchProcessRepo
//}
