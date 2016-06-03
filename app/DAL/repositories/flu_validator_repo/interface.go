package flu_validator_repo

import (
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type IFluValidatorRepo interface {
	GetValidatorsForProject(projectId uuid.UUID, tag string) ([]models.FLUValidator, error)
	Save(*models.FLUValidator) error
}
