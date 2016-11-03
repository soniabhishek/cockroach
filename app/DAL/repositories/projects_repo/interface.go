package projects_repo

import (
	"errors"
	"sync"

	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type IProjectsRepo interface {
	GetById(uuid.UUID) (models.Project, error)
	Add(models.Project) error
	Update(models.Project) error
	Delete(id uuid.UUID) error
	IfIdExist(uuid.UUID) (bool, error)
}

//=============================================================================================//

type inMemProjectRepo struct {
	sync.RWMutex
	projects map[uuid.UUID]models.Project
}

var _ IProjectsRepo = (*inMemProjectRepo)(nil)

func (i *inMemProjectRepo) GetById(id uuid.UUID) (pr models.Project, err error) {
	i.RLock()
	defer i.RUnlock()

	pr, ok := i.projects[id]
	if !ok {
		return pr, errors.New("not found")
	}
	return pr, nil
}

func (i *inMemProjectRepo) Add(models.Project) error {
	return nil
}
func (i *inMemProjectRepo) Update(models.Project) error {
	return nil
}
func (i *inMemProjectRepo) Delete(uuid.UUID) error {
	return nil
}
func (i *inMemProjectRepo) IfIdExist(uuid.UUID) (bool, error) {
	return false, nil
}

func (i *inMemProjectRepo) Save(pr models.Project) error {
	i.Lock()
	defer i.Unlock()

	i.projects[pr.ID] = pr
	return nil
}

func Mock() IProjectsRepo {
	return &inMemProjectRepo{
		projects: make(map[uuid.UUID]models.Project),
	}
}
