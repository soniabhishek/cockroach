package operations

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crowdflux/angel/app/DAL/repositories/clients_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/project_configuration_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/user_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/workflow_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/utilities/clients/models"
	"github.com/crowdflux/angel/utilities/constants"
	"github.com/lib/pq"
)

type Service struct {
}

func (s *Service) CreateClient(cl *utilModels.Client) (status bool, err error) {

	fmt.Println(cl)
	/**-------- Creating User ---------**/
	userRepo := user_repo.New()
	user := models.User{
		Username:  cl.UserName,
		Password:  sql.NullString{cl.Password, true},
		Gender:    sql.NullString{string(cl.Gender), true},
		FirstName: sql.NullString{string(cl.FirstName), true},
		LastName:  sql.NullString{string(cl.LastName), true},
		Locale:    sql.NullString{constants.Empty, true},
		AvatarUrl: sql.NullString{constants.Empty, true},
		Phone:     sql.NullString{string(cl.Phone), true},
	}
	err = userRepo.Add(&user)
	userId := user.ID
	if err != nil {
		return false, err
	}
	plog.Info("Created User. UserId [" + userId.String() + "]")

	/**-------- Creating Client ---------**/
	clientRepo := clients_repo.New()
	client := models.Client{
		UserId:  userId,
		Options: models.JsonF{},
	}
	err = clientRepo.Add(&client)
	clientId := client.ID
	clientSecretId := client.ClientSecretUuid
	if err != nil {
		return false, err
	}
	cl.ClientId = clientId
	cl.ClientSecretId = clientSecretId
	plog.Info("Created Client. ClientId [" + clientId.String() + "] ClientSecretId [" + clientSecretId.String() + "]")

	/**-------- Creating Project ---------**/
	projectRepo := projects_repo.New()
	project := models.Project{
		Label:     cl.ProjectLabel,
		Name:      cl.ProjectName,
		ClientId:  clientId,
		CreatorId: userId,
	}
	err = projectRepo.Add(&project)
	if err != nil {
		return false, err
	}
	projectId := project.ID
	cl.ProjectId = projectId
	plog.Info("Created Project. ProjectId [" + projectId.String() + "]")

	/**-------- Configuring Project ---------**/
	err = project_configuration_repo.New().Add(models.ProjectConfiguration{
		ProjectId:   projectId,
		PostBackUrl: cl.Url,
		Headers:     cl.Header,
		CreatedAt:   pq.NullTime{time.Now(), true},
		UpdatedAt:   pq.NullTime{time.Now(), true},
		Options:     models.JsonF{},
	})
	if err != nil {
		return false, err
	}
	plog.Info("Configured Project. ProjectId [" + projectId.String() + "]")

	/**-------- Creating Workflow ---------**/
	wfId := uuid.NewV4()
	wfr := workflow_repo.New()
	err = wfr.Add(&models.WorkFlow{
		ID:        wfId,
		ProjectId: projectId,
		IsDeleted: sql.NullBool{false, true},
		CreatedAt: pq.NullTime{time.Now(), true},
		UpdatedAt: pq.NullTime{time.Now(), true},
	})
	if err != nil {
		return false, err
	}
	plog.Info("Created Workflow. WorkflowId [" + wfId.String() + "]")
	fmt.Println(cl)
	return true, nil
}
