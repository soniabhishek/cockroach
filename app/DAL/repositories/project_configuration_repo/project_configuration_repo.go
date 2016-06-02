package project_configuration_repo

import (
	"errors"

	"gitlab.com/playment-main/angel/app/DAL/repositories"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
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
	err = fps.db.SelectOne(&v, "SELECT * FROM flu_project_service WHERE project_id = $1", projectId)
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
