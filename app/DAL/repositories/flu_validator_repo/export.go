package flu_validator_repo

import (
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
)

func New() IFluValidatorRepo {
	return &fluValidatorRepo{
		db: postgres.GetPostgresClient(),
	}
}
