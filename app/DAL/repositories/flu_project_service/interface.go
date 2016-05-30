package flu_project_service

import (
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type IFluProjectServiceRepo interface {
	Save(*models.FluProjectService) error
	Get(projectId uuid.UUID) (models.FluProjectService, error)
}
