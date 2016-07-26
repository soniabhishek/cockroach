package flu_validator

import (
	"github.com/crowdflux/angel/app/DAL/repositories/flu_validator_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
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
