package flu_validator_repo

import "gitlab.com/playment-main/angel/app/DAL/clients"

func New() IFluValidatorRepo {
	return &fluValidatorRepo{
		db: clients.GetPostgresClient(),
	}
}
