package feed_line_repo

import (
	"errors"

	"time"

	"fmt"

	"bytes"
	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/DAL/repositories/queries"
	"github.com/crowdflux/angel/app/DAL/repositories/step_repo"
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/lib/pq"
	"strings"
)

type fluRepo struct {
	Db       repositories.IDatabase
	stepRepo step_repo.IStepRepo
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
	flu.CreatedAt = pq.NullTime{time.Now(), true}
	return e.Db.Insert(&flu)
}

func (e *fluRepo) Update(flu models.FeedLineUnit) error {
	flu.UpdatedAt = pq.NullTime{time.Now(), true}
	_, err := e.Db.Update(&flu)
	return err
}

//Insert into Postgress
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

	var queryBuffer bytes.Buffer

	queryBuffer.WriteString(`update feed_line as fl set
		    build = tmp.build, updated_at = tmp.updated_at
		  from (values `)

	updatableRowsCount := 0

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

		if updatableRowsCount > 0 {
			queryBuffer.WriteString(",")
		}
		updatableRowsCount++

		dbFlu.Build.Merge(flus[i].Build)

		idVal, _ := flus[i].ID.Value()
		buildVal, _ := dbFlu.Build.Value()
		updatedAtVal := pq.NullTime{time.Now(), true}.Time.Format(time.RFC3339)

		tmp := fmt.Sprintf(`('%v'::uuid, '%v'::jsonb, '%v'::timestamp with time zone)`, idVal, buildVal, updatedAtVal)
		queryBuffer.WriteString(tmp)

	}
	queryBuffer.WriteString(`) as tmp(id, build, updated_at)
		where tmp.id = fl.id;`)

	if updatableRowsCount == 0 {
		return errors.New("No updatable flu")
	}

	query := queryBuffer.String()
	if config.IsDevelopment() || config.IsStaging() {
		plog.Trace("Running Q: ", query)
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

func (e *fluRepo) BulkFluBuildUpdateByStepType(flus []models.FeedLineUnit, stepType step_type.StepType) (updatedFlus []models.FeedLineUnit, err error) {

	plog.Info("FLUREPO", "Reached Bulk Upload")

	updatableRows, err := e.getUpdableFlus(flus, stepType)

	if err != nil {
		return updatableRows, err
	}

	plog.Info("FLUREPO", "got updatable flus")

	if len(updatableRows) == 0 {
		return updatableRows, ErrNoUpdatableFlus
	}

	var queryBuffer bytes.Buffer

	queryBuffer.WriteString(`update feed_line as fl set
		    build = tmp.build, updated_at = tmp.updated_at
		  from (values `)

	for i, flu := range updatableRows {

		idVal, _ := flu.ID.Value()
		buildVal := flu.Build.String()

		//Replace ' with '' for psql
		buildVal = strings.Replace(buildVal, "'", "''", -1)

		updatedAtVal := pq.NullTime{time.Now(), true}.Time.Format(time.RFC3339)

		if i > 0 {
			queryBuffer.WriteString(",")
		}

		tmp := fmt.Sprintf(`('%v'::uuid, '%v'::jsonb, '%v'::timestamp with time zone)`, idVal, buildVal, updatedAtVal)
		queryBuffer.WriteString(tmp)
	}

	queryBuffer.WriteString(`) as tmp(id, build, updated_at)
		where tmp.id = fl.id;`)

	plog.Info("FLU REPO", "query built")

	res, err := e.Db.Exec(queryBuffer.String())
	if err != nil {
		return updatableRows, err
	}
	if rows, _ := res.RowsAffected(); rows != int64(len(flus)) {
		return updatableRows, ErrPartiallyUpdatedFlus
	}
	return updatableRows, nil
}

func (e *fluRepo) getUpdableFlus(flus []models.FeedLineUnit, stepType step_type.StepType) (updatedFlus []models.FeedLineUnit, err error) {
	type StepTypeMap map[uuid.UUID]step_type.StepType

	var stepTypeMap StepTypeMap = make(StepTypeMap)

	updatableRows := make([]models.FeedLineUnit, 0, len(flus))
	for _, flu := range flus {

		if flu.ID == uuid.Nil {
			//return errors.New("flu not present")
			continue
		}

		dbFlu, err := e.GetById(flu.ID)
		if err != nil {
			plog.Info(err.Error())
			continue
		}
		dbStepType, ok := stepTypeMap[dbFlu.StepId]
		if !ok {
			step, err := e.stepRepo.GetById(dbFlu.StepId)
			if err != nil {
				plog.Error("flurepo", err)
				return []models.FeedLineUnit{}, err
			}
			stepTypeMap[dbFlu.StepId] = step.Type
			dbStepType = step.Type
		}

		if dbStepType != stepType {
			plog.Info("flurepo", "flu doesnot belong to this step")
			continue
		}

		dbFlu.Build.Merge(flu.Build)

		updatableRows = append(updatableRows, dbFlu)

	}

	return updatableRows, nil
}

func (e *fluRepo) GetFlusNotSent(StepId uuid.UUID) (flus []models.FeedLineUnit, err error) {

	//	_, err = e.Db.Select(&flus, `SELECT fl.*
	//FROM feed_line_log fll
	//INNER JOIN feed_line fl on fl.id = fll.flu_id
	//WHERE fll.step_id =$1
	//AND fll.meta_data->>'HttpStatusCode'!='200'
	//and fll.created_at > date_trunc('day',now()-interval '30 days')
	//AND fll.flu_id not in (
	//  SELECT fll.flu_id
	//  FROM feed_line_log fll
	//  WHERE step_id = $2
	//        AND created_at > date_trunc('day',now()-interval '30 days')
	//        AND (fll.meta_data ->> 'HttpStatusCode' = '200') AND fll.message = 'SUCCESS'
	//);`, StepId, StepId)

	_, err = e.Db.Select(&flus, `
SELECT fl.*
FROM feed_line fl LEFT OUTER JOIN
  feed_line_log fll on (fl.id = fll.flu_id
    AND fll.step_id = $1
    AND fll.meta_data->>'HttpStatusCode'='200')
  WHERE fl.step_id = $2
  AND fll.id is NULL;`, StepId, StepId)

	return
}
