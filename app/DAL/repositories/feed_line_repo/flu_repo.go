package feed_line_repo

import (
	"errors"

	"time"

	"fmt"
	"github.com/lib/pq"
	"gitlab.com/playment-main/angel/app/DAL/repositories"
	"gitlab.com/playment-main/angel/app/DAL/repositories/queries"
	"gitlab.com/playment-main/angel/app/config"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/plog"
)

type fluRepo struct {
	Db repositories.IDatabase
}

var _ IFluRepo = &fluRepo{}

//Gets a Flu from database for the given id
func (e *fluRepo) GetById(id uuid.UUID) (models.FeedLineUnit, error) {

	var flu models.FeedLineUnit

	err := e.Db.SelectOne(&flu, "select * from feed_line where id = $1", id)
	if err != nil {
		return flu, err
	}
	return flu, nil
}

func (e *fluRepo) Save(i models.FeedLineUnit) {
	panic(errors.New("Not implemented"))
}

func (e *fluRepo) Add(flu models.FeedLineUnit) error {

	flu.ID = queries.EnsureId(flu.ID)
	return e.Db.Insert(&flu)
}

func (e *fluRepo) Update(flu models.FeedLineUnit) error {
	_, err := e.Db.Update(&flu)
	return err
}

func (e *fluRepo) BulkInsert(flus []models.FeedLineUnit) error {

	var flusInterface []interface{} = make([]interface{}, len(flus))
	for i, _ := range flus {

		if flus[i].ID == uuid.Nil {
			flus[i].ID = uuid.NewV4()
		}
		flus[i].CreatedAt = pq.NullTime{time.Now(), true}
		flusInterface[i] = &flus[i]
	}

	err := e.Db.Insert(flusInterface...)
	return err
}

//Gets a Flu from database for the given id
func (e *fluRepo) GetByStepId(StepId uuid.UUID) ([]models.FeedLineUnit, error) {

	fmt.Println(StepId)
	var flus []models.FeedLineUnit
	_, err := e.Db.Select(&flus, "select * from feed_line where step_id = $1", StepId)
	if err != nil {
		return flus, err
	}
	return flus, nil
}

func (e *fluRepo) BulkUpdate(flus []models.FeedLineUnit) error {

	var flusInterface []interface{} = make([]interface{}, len(flus))
	for i, _ := range flus {

		if flus[i].ID == uuid.Nil {
			return errors.New("flu not present")
		}
		flus[i].UpdatedAt = pq.NullTime{time.Now(), true}
		flusInterface[i] = &flus[i]
	}

	total, err := e.Db.Update(flusInterface...)
	if total != int64(len(flus)) {
		err = errors.New("Partially dumped the data. [" + err.Error() + "]")
	}
	return err
}

func (e *fluRepo) BulkFluBuildUpdate(flus []models.FeedLineUnit) error {
	query := `update feed_line as fl set
		    build = tmp.build, updated_at = tmp.updated_at
		  from (values `

	l := len(flus)
	for i, _ := range flus {

		if flus[i].ID == uuid.Nil {
			//return errors.New("flu not present")
			continue
		}

		dbFlu, err := e.GetById(flus[i].ID)
		if err != nil {
			plog.Info(err.Error())
			continue
		}
		flus[i].ReferenceId = dbFlu.ReferenceId
		flus[i].Data = dbFlu.Data
		flus[i].Tag = dbFlu.Tag
		flus[i].CreatedAt = dbFlu.CreatedAt
		flus[i].StepId = dbFlu.StepId
		flus[i].ProjectId = dbFlu.ProjectId

		idVal, _ := flus[i].ID.Value()
		buildVal, _ := flus[i].Build.Value()
		updatedAtVal := pq.NullTime{time.Now(), true}.Time.Format(time.RFC3339)

		tmp := fmt.Sprintf(`('%v'::uuid, '%v'::jsonb, '%v'::timestamp with time zone)`, idVal, buildVal, updatedAtVal)
		query += tmp
		if i < l-1 {
			query += ","
		}
	}
	query += `) as tmp(id, build, updated_at)
		where tmp.id = fl.id;`

	if config.IsDevelopment() || config.IsStaging() {
		plog.Info("Running Q: ", query)
	}
	res, err := e.Db.Exec(query)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows != int64(len(flus)) {
		err = errors.New("Partially dumped the data.")
	}
	return err
}
