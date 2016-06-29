package user_repo

import (
	"gitlab.com/playment-main/angel/app/DAL/clients/postgres"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
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
