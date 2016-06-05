package projects_repo

import (
	"errors"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"sync"
)

type IProjectsRepo interface {
	GetById(uuid.UUID) (models.Project, error)
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
