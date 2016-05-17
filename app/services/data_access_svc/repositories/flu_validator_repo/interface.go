package flu_validator_repo

import (
	"gitlab.com/playment-main/support/app/models"
	"gitlab.com/playment-main/support/app/models/uuid"
)

type IFluValidatorRepo interface {
	GetValidatorsForMacroTask(macroTaskId uuid.UUID, tag string) ([]models.FLUValidator, error)
	Save(*models.FLUValidator) error
}
