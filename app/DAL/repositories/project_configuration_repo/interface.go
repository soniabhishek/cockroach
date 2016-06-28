package project_configuration_repo

import (
	"errors"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"sync"
)

type IProjectConfigurationRepo interface {
	Save(*models.ProjectConfiguration) error
	Get(projectId uuid.UUID) (models.ProjectConfiguration, error)
	Add(models.ProjectConfiguration) error
	Update(models.ProjectConfiguration) error
	Delete(uuid.UUID) error
}

//=============================================================================================//

type inMemProjectConfigurationRepo struct {
	sync.RWMutex
	projectConfigs map[uuid.UUID]models.ProjectConfiguration
}

var _ IProjectConfigurationRepo = (*inMemProjectConfigurationRepo)(nil)

func (i *inMemProjectConfigurationRepo) Get(projectId uuid.UUID) (pr models.ProjectConfiguration, err error) {
	i.RLock()
	defer i.RUnlock()

	pr, ok := i.projectConfigs[projectId]
	if !ok {
		return pr, errors.New("not found")
	}
	return pr, nil
}

func (i *inMemProjectConfigurationRepo) Save(pr *models.ProjectConfiguration) error {
	i.Lock()
	defer i.Unlock()

	i.projectConfigs[pr.ProjectId] = *pr
	return nil
}

func Mock() IProjectConfigurationRepo {
	return &inMemProjectConfigurationRepo{
		projectConfigs: make(map[uuid.UUID]models.ProjectConfiguration),
	}
}

func (i *inMemProjectConfigurationRepo) Add(pr models.ProjectConfiguration) error {
	return nil
}

func (i *inMemProjectConfigurationRepo) Update(pr models.ProjectConfiguration) error {
	return nil
}

func (i *inMemProjectConfigurationRepo) Delete(id uuid.UUID) error {
	return nil
}
