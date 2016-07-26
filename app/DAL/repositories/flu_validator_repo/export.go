package flu_validator_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
)

func New() IFluValidatorRepo {
	return &fluValidatorRepo{
		db: postgres.GetPostgresClient(),
	}
}
