package user_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/clients_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/user_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IUserService interface {
	CreateUser(models.User) error
	GetUser(uuid.UUID) (models.User, error)
}

func New() IUserService {
	return &userService{
		userRepo: user_repo.New(),
	}
}
