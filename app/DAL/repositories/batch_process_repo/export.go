package batch_process_repo

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IBatchProcessRepo interface {
	Get(uuid.UUID) (models.BatchProcess, error)
	Save(models.BatchProcess) error
}

//func New() IBatchProcessRepo {
//	return IBatchProcessRepo
//}
