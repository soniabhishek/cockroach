package flu_project_service

import (
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/DAL/repositories"
	"errors"
)

const fluProjectServiceTable string = "flu_project_service"



type fluProjectService struct {
	db repositories.IDatabase
}

var _ IFluProjectServiceRepo = &fluProjectService{}

//Unused
func (fps *fluProjectService) Save(fv *models.FluProjectService) error {
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

func (fps *fluProjectService) Get(projectId uuid.UUID) (models.FluProjectService, error){
	return fps.get(projectId)
}

/*-------------------------------------------------------------------------------------------------------------*/
/*----------------------------------Basic DB Queries Implementation--------------------------------------------*/
/*-------------------------------------------------------------------------------------------------------------*/
func (fps *fluProjectService) get(projectId uuid.UUID) (v models.FluProjectService, err error) {
	err = fps.db.SelectOne(&v, "SELECT * FROM flu_project_service WHERE project_id = $1", projectId)
	return
}

func (fps *fluProjectService) update(fp models.FluProjectService) error {

	_, err := fps.db.Update(&fp)
	return err
}

func (fps *fluProjectService) insertMany(projects []models.FluProjectService) error {

	err := fps.db.Insert(toInterfaceArray(projects)...)
	return err
}

func (fps *fluProjectService) deleteMany(projects []models.FluProjectService) error {

	_, err := fps.db.Delete(toInterfaceArray(projects)...)
	return err
}

func toInterfaceArray(fps []models.FluProjectService) []interface{} {
	var projects []interface{} = make([]interface{}, len(fps))
	for i, _ := range fps {
		projects[i] = &fps[i]
	}
	return projects
}
