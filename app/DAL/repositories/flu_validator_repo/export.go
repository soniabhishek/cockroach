package flu_validator_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/DAL/imdb"
)

func New() IFluValidatorRepo {
	return &fluValidatorRepo{
		db: postgres.GetPostgresClient(),
		imdb: imdb.NewFluValidateCache(),
	}
}
