package user_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/user_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/lib/pq"
	"time"
)

type userService struct {
	userRepo user_repo.IUserRepo
}

func (cs *userService) CreateUser(user models.User) error {
	user.ID = uuid.NewV4()
	currentTime := time.Now()
	user.CreatedAt = pq.NullTime{currentTime, true}
	user.UpdatedAt = pq.NullTime{currentTime, true}
	err := cs.userRepo.Add(user)
	return err
}

func (cs *userService) GetUser(userId uuid.UUID) (user models.User, err error) {
	return
}
