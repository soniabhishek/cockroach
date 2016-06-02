package project_configuration_repo

import (
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type IProjectConfigurationRepo interface {
	Save(*models.ProjectConfiguration) error
	Get(projectId uuid.UUID) (models.ProjectConfiguration, error)
}
