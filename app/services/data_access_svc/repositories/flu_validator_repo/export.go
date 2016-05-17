package flu_validator_repo

import "gitlab.com/playment-main/support/app/services/data_access_svc/clients"

func New() IFluValidatorRepo {
	return &fluValidatorRepo{
		db: clients.GetPostgresClient(),
	}
}
