package flu_validator

import (
	"gitlab.com/playment-main/angel/app/DAL/repositories/flu_validator_repo"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

func New() IFluValidatorService {
	return &fluValidator{
		fluValidatorRepo: flu_validator_repo.New(),
	}
}

type IFluValidatorService interface {
	GetValidators(projectId uuid.UUID, tag string) (fvs []models.FLUValidator, err error)
	SaveValidator(fv *models.FLUValidator) (err error)
	Validate(flu models.FeedLineUnit) (IsValid bool, err error)
}
