package project_configuration_repo

import (
	"errors"

	"fmt"

	"github.com/crowdflux/angel/app/DAL/repositories"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

const projectConfigurationTable string = "project_configuration"

type projectConfigurationRepo struct {
	db repositories.IDatabase
}

var _ IProjectConfigurationRepo = &projectConfigurationRepo{}

//Unused
func (fps *projectConfigurationRepo) Save(fv *models.ProjectConfiguration) error {
	/*existingFv, err := f.get(fv.ID)

	var exists bool

	if err == nil {
		exists = true
	}

	if exists {

		//It requires to copy the existing createdAt otherwise it will update it as nil
		fv.CreatedAt = existingFv.CreatedAt
		fv.UpdatedAt = gorp.NullTime{time.Now(), true}
		_, err = f.db.Update(fv)
	} else {

		fv.CreatedAt = gorp.NullTime{time.Now(), true}
		if fv.ID == uuid.Nil {
			fv.ID = uuid.NewV4()
		}
		err = f.db.Insert(fv)
	}
	return err*/
	return errors.New("UnUsed/Uninitialized method")
}

func (fps *projectConfigurationRepo) Get(projectId uuid.UUID) (models.ProjectConfiguration, error) {
	return fps.get(projectId)
}

/*-------------------------------------------------------------------------------------------------------------*/
/*----------------------------------Basic DB Queries Implementation--------------------------------------------*/
/*-------------------------------------------------------------------------------------------------------------*/
func (fps *projectConfigurationRepo) get(projectId uuid.UUID) (v models.ProjectConfiguration, err error) {
	err = fps.db.SelectOne(&v, "SELECT * FROM "+projectConfigurationTable+" WHERE project_id = $1", projectId)
	return
}

func (fps *projectConfigurationRepo) update(fp models.ProjectConfiguration) error {

	_, err := fps.db.Update(&fp)
	return err
}

func (fps *projectConfigurationRepo) insertMany(projects []models.ProjectConfiguration) error {

	err := fps.db.Insert(toInterfaceArray(projects)...)
	return err
}

func (fps *projectConfigurationRepo) deleteMany(projects []models.ProjectConfiguration) error {

	_, err := fps.db.Delete(toInterfaceArray(projects)...)
	return err
}

func toInterfaceArray(fps []models.ProjectConfiguration) []interface{} {
	var projects []interface{} = make([]interface{}, len(fps))
	for i, _ := range fps {
		projects[i] = &fps[i]
	}
	return projects
}

func (fps *projectConfigurationRepo) Add(pj models.ProjectConfiguration) error {
	return fps.db.Insert(&pj)
}

func (fps *projectConfigurationRepo) Update(pj models.ProjectConfiguration) error {
	_, err := fps.db.Update(&pj)
	return err
}

func (fps *projectConfigurationRepo) Delete(id uuid.UUID) error {
	query := fmt.Sprintf(`delete from project_configuration where project_id='%v'::uuid`, id)
	res, err := fps.db.Exec(query)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows < 1 {
		err = errors.New("Could not delete ProjectConfiguration with ID [" + id.String() + "]")
	}
	return err
}
