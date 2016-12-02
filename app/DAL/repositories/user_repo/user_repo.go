package user_repo

import (
	"errors"
	"fmt"

	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/lib/pq"
	"time"
)

type user_repo struct {
	db repositories.IDatabase
}

var _ IUserRepo = &user_repo{}

func (ur *user_repo) Add(u *models.User) error {
	u.ID = uuid.NewV4()
	u.UpdatedAt = pq.NullTime{time.Now(), true}
	u.CreatedAt = u.UpdatedAt
	return ur.db.Insert(u)
}

func (ur *user_repo) Update(u models.User) error {
	_, err := ur.db.Update(&u)
	return err
}

func (ur *user_repo) Delete(id uuid.UUID) error {
	query := fmt.Sprintf(`delete from users where id='%v'::uuid`, id)
	res, err := ur.db.Exec(query)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows < 1 {
		err = errors.New("Could not delete user with ID [" + id.String() + "]")
	}
	return err
}

func (ur *user_repo) IfIdExist(id uuid.UUID) (ifExist bool, err error) {
	err = ur.db.SelectOne(&ifExist, `select exists(select 1 from users where id=$1)`, id)
	return
}
