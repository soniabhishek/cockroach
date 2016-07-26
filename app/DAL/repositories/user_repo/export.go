package user_repo

import (
	"github.com/crowdflux/angel/app/DAL/clients/postgres"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IUserRepo interface {
	Add(models.User) error
	Update(models.User) error
	Delete(uuid.UUID) error
}

func New() IUserRepo {
	return &user_repo{
		db: postgres.GetPostgresClient(),
	}
}
