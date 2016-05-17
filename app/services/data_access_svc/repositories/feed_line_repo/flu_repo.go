package feed_line_repo

import (
	"errors"

	"time"

	"gitlab.com/playment-main/support/app/models"
	"gitlab.com/playment-main/support/app/models/uuid"
	"gitlab.com/playment-main/support/app/services/data_access_svc/repositories"
	"gopkg.in/gorp.v1"
)

type fluRepo struct {
	Db repositories.IDatabase
}

var _ IFluRepo = &fluRepo{}

//Gets a Flu from database for the given id
func (e *fluRepo) GetById(id uuid.UUID) (models.FeedLineUnit, error) {

	var ip models.FeedLineUnit

	err := e.Db.SelectOne(&ip, "select * from feedline where id = $1", id)
	if err != nil {
		return ip, err
	}
	return ip, nil
}

func (e *fluRepo) Save(i models.FeedLineUnit) {
	panic(errors.New("Not implemented"))
}

func (e *fluRepo) BulkInsert(flus []models.FeedLineUnit) error {

	var flusInterface []interface{} = make([]interface{}, len(flus))
	for i, _ := range flus {

		if flus[i].ID == uuid.Nil {
			flus[i].ID = uuid.NewV4()
		}
		flus[i].CreatedAt = gorp.NullTime{time.Now(), true}
		flusInterface[i] = &flus[i]
	}

	err := e.Db.Insert(flusInterface...)
	return err
}
