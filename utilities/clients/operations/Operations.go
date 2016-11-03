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
	userId := uuid.NewV4()
	err = userRepo.Add(models.User{
		ID:                      userId,
		Username:                cl.UserName,
		Password:                sql.NullString{cl.Password, true},
		CreatedAt:               pq.NullTime{time.Now(), true},
		UpdatedAt:               pq.NullTime{time.Now(), true},
		Gender:                  sql.NullString{string(cl.Gender), true},
		FirstName:               sql.NullString{string(cl.FirstName), true},
		LastName:                sql.NullString{string(cl.LastName), true},
		Locale:                  sql.NullString{constants.Empty, true},
		AvatarUrl:               sql.NullString{constants.Empty, true},
		IncorrectQuestionsCount: 0,
		CorrectQuestionsCount:   0,
		PendingQuestionsCount:   0,
		CoinsCount:              0,
		CurrentPower:            0,
		CouponRedeemedCount:     0,
		Phone:                   sql.NullString{string(cl.Phone), true},
		TotalCoinsCount:         0,
	})
	if err != nil {
		return false, err
	}
	plog.Info("Created User. UserId [" + userId.String() + "]")

	/**-------- Creating Client ---------**/
	clientRepo := clients_repo.New()
	clientId, clientSecretId := uuid.NewV4(), uuid.NewV4()
	err = clientRepo.Add(models.Client{
		ID:               clientId,
		UserId:           userId,
		ClientSecretUuid: clientSecretId,
		CreatedAt:        pq.NullTime{time.Now(), true},
		UpdatedAt:        pq.NullTime{time.Now(), true},
		Options:          models.JsonF{},
	})
	if err != nil {
		return false, err
	}
	cl.ClientId = clientId
	cl.ClientSecretId = clientSecretId
	plog.Info("Created Client. ClientId [" + clientId.String() + "] ClientSecretId [" + clientSecretId.String() + "]")

	/**-------- Creating Project ---------**/
	projectId := uuid.NewV4()
	projectRepo := projects_repo.New()
	err = projectRepo.Add(models.Project{
		ID:        projectId,
		Label:     cl.ProjectLabel,
		Name:      cl.ProjectName,
		ClientId:  clientId,
		CreatorId: userId,
		StartedAt: pq.NullTime{time.Now(), true},
		EndedAt:   pq.NullTime{time.Now(), false},
		CreatedAt: pq.NullTime{time.Now(), true},
		UpdatedAt: pq.NullTime{time.Now(), true},
	})
	if err != nil {
		return false, err
	}
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
